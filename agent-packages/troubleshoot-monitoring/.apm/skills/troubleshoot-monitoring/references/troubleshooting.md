# Qubership Monitoring Operator — troubleshooting

This reference covers failures of an already-deployed monitoring stack — problems that surface while installing the
release into a cluster or while operating it afterward. It does not cover building the operator or component images, CI,
or local development.

Cases marked `Derived from source` are compiled from this repository's operator code, Helm values, and documentation.
Cases marked as compiled from a sibling repository come from the code of a component that qubership-monitoring-operator
deploys and are not all yet confirmed on a live installation. Cases under components sourced only from upstream issues
and vendor documentation are marked as compiled from external research and are likewise pending confirmation on a live
installation. Verify version-specific and code-derived cases against the build you actually run before acting on a
destructive step.

## Deploy and CRDs

### Rendered manifests contain a new resource that already exists

**Symptoms:**

* A deploy job or a manual `helm install` / `helm upgrade` fails with a resource conflict:

  <!-- markdownlint-disable line-length -->
  ```text
  helm.go:75: [debug] existing resource conflict: kind: <kind>, namespace: <namespace>, name: <name>
  rendered manifests contain a new resource that already exists. Unable to continue with [install|update]
  ```
  <!-- markdownlint-enable line-length -->

* The conflicting object can be any managed kind (for example `PlatformMonitoring`, `Role`, `ClusterRole`).

**Root cause:**

An object of that kind and name already exists in the namespace but was not created by this Helm release, so it carries
no Helm ownership metadata. Helm refuses to adopt or overwrite an object it does not track.

**How to check:**

1. Read the object's Helm release annotation. Empty output means Helm does not own it:

   ```bash
   kubectl -n <namespace> get <kind> <name> -o jsonpath='{.metadata.annotations.meta\.helm\.sh/release-name}{"\n"}'
   ```

2. A Helm-owned object shows both `meta.helm.sh/release-name` and `meta.helm.sh/release-namespace` annotations and the
   `app.kubernetes.io/managed-by: Helm` label. If any are missing, the object is untracked.

**How to fix:**

1. Preferred, non-destructive: adopt the existing object into the release by adding the two ownership annotations and
   the managed-by label. Substitute the release name and namespace your deploy uses:

   ```bash
   kubectl -n <namespace> annotate <kind> <name> \
     meta.helm.sh/release-name=<release-name> \
     meta.helm.sh/release-namespace=<release-namespace> --overwrite
   kubectl -n <namespace> label <kind> <name> app.kubernetes.io/managed-by=Helm --overwrite
   ```

   Re-run the deploy. Helm now treats the object as its own.

2. **DANGEROUS — deletes the pre-existing object and any data or configuration it holds; a `PlatformMonitoring` CR
   deletion tears down the managed stack.** Only if the object is disposable, delete it and let the deploy recreate it:

   ```bash
   kubectl -n <namespace> delete <kind> <name>
   ```

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### PlatformMonitoring or component CRD apply fails on annotation or etcd size limit

**Symptoms:**

* `kubectl apply -f <crd>` for `platformmonitorings.monitoring.netcracker.com` (about 5 MB) fails, or the object is
  rejected by the API server for exceeding the per-object size limit.
* The failure references the `kubectl.kubernetes.io/last-applied-configuration` annotation or an etcd request size.

**Root cause:**

`kubectl apply` writes the full object into the `last-applied-configuration` annotation, roughly doubling the stored
size. The `PlatformMonitoring` CRD is already about 5 MB, so the annotated copy exceeds the etcd per-resource limit. See
also **GrafanaDashboard rejected because it exceeds the etcd 1 MB limit**, the same annotation-and-size mechanism for
dashboard objects.

**How to check:**

1. Confirm the CRD size and that a client-side apply is being used. Inspect the manifest or the deploy command for
   `kubectl apply` on the CRDs.
2. If the object already exists, check whether it carries the bloating annotation:

   ```bash
   kubectl get crd platformmonitorings.monitoring.netcracker.com \
     -o jsonpath='{.metadata.annotations.kubectl\.kubernetes\.io/last-applied-configuration}{"\n"}' | wc -c
   ```

**How to fix:**

1. Install CRDs with server-side apply, which does not write the `last-applied-configuration` annotation:

   ```bash
   kubectl apply -f charts/qubership-monitoring-crds/crds/ --server-side
   ```

2. Alternatively, create or replace CRDs without apply:

   ```bash
   kubectl create -f <crds-directory>/
   kubectl replace -f <crds-directory>/
   ```

   Never use `kubectl apply` for these CRDs.

**How to avoid this issue:**

Install the standalone `qubership-monitoring-crds` chart, or use ArgoCD with `ServerSideApply=true`, so large CRDs never
receive the `last-applied-configuration` annotation.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `README.md`, `docs/user-guides/manual-create-crds.md`
<!-- markdownlint-enable line-length -->

## Monitoring operator

### Prometheus or Grafana operators are restarting

**Symptoms:**

* `prometheus-operator` or `grafana-operator` restarts shortly after starting.
* Pod status or events contain `OOMKilled`.
* Operator memory usage sits near the configured memory limit.

**Root cause:**

Both operators allocate a burst of memory at startup to process every Custom Resource they manage. On clusters with many
`ServiceMonitor` and `GrafanaDashboard` CRs, that startup spike can exceed the memory limit and the pod is killed before
it settles. Observed peaks were above `400Mi` for grafana-operator and above `200Mi` for prometheus-operator on a
cluster with 835 `ServiceMonitor` and 347 `GrafanaDashboard` CRs, dropping to roughly `143M` and `49Mi` after
processing. OpenShift `v4.x` needs more memory for prometheus-operator than earlier versions or Kubernetes with the same
configuration, so it hits this more often.

**How to check:**

1. Inspect the operator pod's restart count and last state:

   ```bash
   kubectl -n <namespace> get pod <operator-pod> -o jsonpath='{.status.containerStatuses[*].lastState}{"\n"}'
   ```

2. Look for `OOMKilled` in the pod description events, and compare live memory usage against the limit:

   ```bash
   kubectl -n <namespace> describe pod <operator-pod>
   kubectl -n <namespace> top pod <operator-pod>
   ```

**How to fix:**

1. Raise the operator memory limit and redeploy, or edit the limits directly in the `PlatformMonitoring` CR. Increase in
   `100Mi` steps until the restarts stop:

   ```yaml
   grafana:
     operator:
       resources:
         requests:
           cpu: 50m
           memory: 256Mi
         limits:
           cpu: 100m
           memory: 512Mi
   prometheus:
     operator:
       resources:
         requests:
           cpu: 100m
           memory: 256Mi
         limits:
           cpu: 200m
           memory: 512Mi
   ```

   The memory `limit` must stay at or above the `requests` value and be raised, not lowered — a limit below the request
   is rejected by the API server, and a limit under the operator's real usage re-triggers the same OOMKill.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### Unable to update etcd certificates due to a lack of permission

**Symptoms:**

* etcd metrics stop updating because the etcd client certificates are never copied into the monitoring Secret.
* The `etcd-certs-to-secret` Job (or its CronJob) fails, logging an error while reading the etcd certificates or writing
  the Secret:

  <!-- markdownlint-disable line-length -->
  ```text
  {"level":"ERROR","msg":"etcd certificates synchronization failed","error":"..."}
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

The etcd client certificates are copied into the monitoring namespace by the `etcd-certs-to-secret` Job
(`cmd/etcd-certs-to-secret/main.go`), not by the operator itself. The Job runs under its own ServiceAccount
(`.Values.etcdCertsJob.rbac.serviceAccountName`, default `etcd-certs-to-secret-sa`) bound to a Role
(`.Values.etcdCertsJob.rbac.roleName`) and a ClusterRole (`.Values.etcdCertsJob.rbac.clusterRoleName`). When that
ServiceAccount lacks read access to the source etcd secrets and configmaps (`etcd-client`, `etcd-metric-client`,
`etcd-serving-ca`, `etcd-metric-serving-ca`) or write access to the target Secret, the copy fails and etcd metrics go
stale. The older operator-emitted WARN line about "a lack of permission to access the requested etcd resource" is not
present in the current operator Go source; attribute the failure to the Job. See also **Cluster-scope operations are
forbidden with privileged rights disabled**.

**How to check:**

1. Read the `etcd-certs-to-secret` Job or CronJob pod log for the synchronization error:

   ```bash
   kubectl -n <monitoring-namespace> logs job/<etcd-certs-to-secret-job>
   ```

2. Confirm the Job's ServiceAccount can read the source etcd certificates and write the target Secret:

   <!-- markdownlint-disable line-length -->
   ```bash
   kubectl auth can-i get secrets \
     --as=system:serviceaccount:<monitoring-namespace>:etcd-certs-to-secret-sa -n <monitoring-namespace>
   ```
   <!-- markdownlint-enable line-length -->

   A `no` result confirms the RBAC gap.

**How to fix:**

1. Grant the etcd RBAC to the Job's ServiceAccount, not the operator ClusterRole. Confirm the Role
   (`.Values.etcdCertsJob.rbac.roleName`) and ClusterRole (`.Values.etcdCertsJob.rbac.clusterRoleName`) the chart
   renders are applied and bound to `.Values.etcdCertsJob.rbac.serviceAccountName`, so the Job can read the source etcd
   secrets and configmaps and write the monitoring Secret. Re-run the Job (or wait for the CronJob) once the RBAC is in
   place.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `cmd/etcd-certs-to-secret/main.go`, `charts/qubership-monitoring-operator/templates/etcd-certs-to-secret-job/role.yaml`
<!-- markdownlint-enable line-length -->

### A component was removed after its Install flag was set to false

**Symptoms:**

* After a redeploy, a previously working component (Grafana, Prometheus, an exporter) and its workloads are gone.
* The operator log for that component reads `Uninstalling component if exists`.
* Dashboards, alerts, or scraped metrics that the component served stop.

**Root cause:**

Each sub-reconciler runs `uninstall(cr)` whenever the component's `Install` flag is not true. Setting `Install: false`
(or removing the sub-spec) on an already-deployed component does not merely stop managing it — it deletes the running
workloads the operator created (Deployment, Service, Ingress, ServiceMonitor, and config). PersistentVolumeClaims are
generally left in place: `PushgatewayReconciler.uninstall`, for example, deletes the Deployment, Service, Ingress, and
routes but not the PVC. Do not rely on either behavior — treat neither the removal of the workloads nor the retention of
a PVC as guaranteed for a given component.

**How to check:**

1. Compare the current CR against the previous state for the component's `install` flag:

   ```bash
   kubectl -n <namespace> get platformmonitoring <cr-name> \
     -o jsonpath='{.spec.grafana.install}{"\n"}'
   ```

2. Search the operator log for `Uninstalling component if exists` around the time the workloads disappeared.

**How to fix:**

1. **DANGEROUS — the uninstall has already removed the component's running workloads and config; a retained PVC may
   still hold the previous data, but recreating a component with a fresh `volumeClaimTemplate` provisions an empty
   volume and does not reattach old data on its own.** Set the component's `Install` back to `true` and redeploy. If a
   PVC survived, confirm the recreated workload binds it; otherwise restore from a backup.

**How to avoid this issue:**

To stop managing a component without deleting it, do not flip `Install` to `false`. Leave it installed and set `Paused:
true` (see the paused-component case below).

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `controllers/grafana/reconciler.go`, `api/v1/platformmonitoring_types.go`
<!-- markdownlint-enable line-length -->

### A change to a component is never applied because it is paused

**Symptoms:**

* Edits to a component's spec in the `PlatformMonitoring` CR have no effect on the running workload.
* The operator log for that component reads `Reconciling paused` and `Component NOT reconciled`.

**Root cause:**

Each sub-reconciler early-returns when the component's `Paused: true` flag is set. The operator stops creating and
updating that component's objects, so any configuration change is silently ignored until the flag is cleared.

**How to check:**

1. Read the component's `paused` flag:

   ```bash
   kubectl -n <namespace> get platformmonitoring <cr-name> \
     -o jsonpath='{.spec.grafana.paused}{"\n"}'
   ```

2. Confirm with the operator log lines `Reconciling paused` / `Component NOT reconciled` for that component.

**How to fix:**

1. Set the component's `Paused` flag back to `false` and redeploy or edit the CR. The next reconcile applies the pending
   changes.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `controllers/grafana/reconciler.go`, `api/v1/platformmonitoring_types.go`
<!-- markdownlint-enable line-length -->

### Cluster-scope operations are forbidden with privileged rights disabled

**Symptoms:**

* With `global.privilegedRights: false`, the operator or a managed operator logs cluster-scope RBAC errors, for example
  cannot list `clusterrolebindings` at the cluster scope, or the etcd certificate WARN described above.
* Cluster-scoped resources the stack expects are never created.

**Root cause:**

When `global.privilegedRights` is `false`, the chart installs only a namespaced `Role` and no `ClusterRole` or
`ClusterRoleBinding`. Any component that needs cluster-scoped access is then forbidden unless the required ClusterRoles
were pre-created out of band. See also **Unable to update etcd certificates due to a lack of permission**, a related
restricted-RBAC failure in the etcd certificate Job.

**How to check:**

1. Read the effective value in the deploy values or on the release:

   ```bash
   grep -n "privilegedRights" charts/qubership-monitoring-operator/values.yaml
   ```

2. List what RBAC the release actually created and confirm no ClusterRole/ClusterRoleBinding exists for the operator:

   ```bash
   kubectl get clusterrole,clusterrolebinding | grep monitoring
   ```

**How to fix:**

1. If the environment permits cluster-scoped RBAC, set `global.privilegedRights: true` and redeploy.
2. If cluster-scoped RBAC must stay disabled, pre-create the required ClusterRoles and ClusterRoleBindings out of band
   and confirm the ServiceAccounts are bound to them before redeploying.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `charts/qubership-monitoring-operator/values.yaml`
<!-- markdownlint-enable line-length -->

## Metric collection

### Metrics absent or errors during metrics collection

**Symptoms:**

* Metrics from a service are missing in Prometheus or VictoriaMetrics even though the pod is running and a
  `ServiceMonitor` / `PodMonitor` exists.
* A target appears in the `Dropped` state, or an active target shows a `Last error`.

**Root cause:**

There are several common causes, ordered by how often they turn out to be the real one:

1. The `ServiceMonitor` / `PodMonitor` is missing the label the operator's selector requires (by convention
   `app.kubernetes.io/component: monitoring`). By default the operator sets an empty selector
   (`serviceMonitorSelector: {}` / `serviceScrapeSelector: {}`), which selects every monitor regardless of labels, so
   the label matters only when the deployment configured a non-empty label selector through the CR.
2. The label selector in the monitor does not match the labels on the Pod and Service, so targets are dropped.
3. Namespace selectors (`serviceMonitorNamespaceSelector` and similar) exclude the namespace.
4. The scrape itself fails: wrong port or path, an endpoint that is unreachable, or a scrape that does not finish within
   `scrapeTimeout`. When the observed scrape duration approaches the configured timeout, scrapes start failing and the
   target shows as down.
5. The configuration has simply not applied yet (it can take 30 seconds to 3 minutes).

**How to check:**

1. Check the scrape configuration. For VictoriaMetrics open the VMAuth or VMAgent Ingress and go to `/config`; for
   Prometheus open its Ingress and go to `Status -> Configuration`. Search for the monitor name and read its relabel
   rules, which encode the filter criteria, for example:

   ```yaml
   relabel_configs:
   - action: keep
     source_labels: [__meta_kubernetes_service_label_component]
     regex: backup-daemon
   ```

2. Check service discovery. For VictoriaMetrics go to `/service-discovery`; for Prometheus go to `Status ->
   ServiceDiscovery`. Find the job and confirm it has at least one target. If every target is `Dropped`, the label
   selectors are wrong.
3. Check the target status. For VictoriaMetrics go to `/targets`; for Prometheus go to `Status -> Targets`. If a target
   is `UP` but metrics are still missing, read the `Last error` column and, for VictoriaMetrics, the `response` link.
4. If a target shows down with a scrape error, read its scrape duration on the targets page and compare it to the
   configured `scrapeTimeout`. A duration near the timeout confirms a scrape-timeout failure rather than a selector or
   discovery problem.

**How to fix:**

1. If the deployment set a non-empty selector, add the label it requires (by convention
   `app.kubernetes.io/component: monitoring`) to the `ServiceMonitor` / `PodMonitor`.
2. Align the monitor's label selector with the labels actually present on the Pod and Service.
3. Adjust the namespace selector so the target namespace is included.
4. For scrape failures, correct the port or path, raise `scrapeTimeout` on the `ServiceMonitor` / `PodMonitor` (keeping
   it below the scrape interval) when the scrape is timing out, or fix the application endpoint.
5. If nothing is wrong, wait up to a few minutes for the configuration to apply, and check victoriametrics-operator /
   prometheus-operator logs.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### VMSingle cannot handle container_start_time_seconds with too small a timestamp

**Symptoms:**

* VMSingle logs a warning about rows with a timestamp outside the retention window:

  <!-- markdownlint-disable line-length -->
  ```text
  2024-01-08T09:42:33.094Z    warn    VictoriaMetrics/lib/storage/storage.go:1733    warn occurred during rows addition:
  cannot insert row with too small timestamp 1696452826250 outside the retention; minimum allowed timestamp is
  1704447752000; probably you need updating -retentionPeriod command-line flag; metricName:
  container_start_time_seconds{...}
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

The kubelet ServiceMonitor scrapes three endpoints (`/metrics/resource`, `/metrics/cadvisor`, and the default path), and
`container_start_time_seconds` appears on several of them. Only the copy from `/metrics/resource` carries a timestamp
old enough to fall outside retention. This is fixable on the Kubernetes side since v1.29. The default monitoring
configuration already drops the metric from that one endpoint, so you only hit this after overriding the default
Kubernetes monitors configuration.

**How to check:**

1. Read the VMSingle log for the `too small timestamp` warning naming `container_start_time_seconds`. This line is
   emitted by VictoriaMetrics only (`lib/storage/storage.go`, tied to `-retentionPeriod`); Prometheus does not emit it,
   so do not search Prometheus for the same wording.
2. Confirm your kubelet ServiceMonitor no longer contains the default drop rules for that metric from
   `/metrics/resource`.

**How to fix:**

1. Restore the relabeling and metric-relabeling rules that drop `container_start_time_seconds` only from the
   `/metrics/resource` endpoint in the `kubernetesMonitors` section:

   ```yaml
   kubernetesMonitors:
     kubeletServiceMonitor:
       metricRelabelings:
         - action: drop
           regex: container_start_time_seconds;\/metrics\/resource
           separator: ;
           sourceLabels: ['__name__', 'metrics_path']
         - action: labeldrop
           regex: metrics_path
       relabelings:
         - action: replace
           regex: (\/metrics\/resource)
           replacement: $1
           sourceLabels:
             - __metrics_path__
           targetLabel: metrics_path
   ```

2. If the ServiceMonitor is not managed by monitoring-operator, you may add the rules directly to it; a managed monitor
   would have such manual edits overwritten on the next reconcile.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

## VictoriaMetrics stack

### VictoriaMetrics operator cannot remove ClusterRole or ClusterRoleBinding

**Symptoms:**

* A deploy in `Clean Install` mode failed on a cluster that already ran the VictoriaMetrics stack.
* The victoriametrics-operator logs an RBAC error:

  <!-- markdownlint-disable line-length -->
  ```text
  Failed to watch *v1.ClusterRoleBinding: failed to list *v1.ClusterRoleBinding: clusterrolebindings.rbac.authorization.k8s.io is forbidden:
    User "system:serviceaccount:system-monitor:system-monitor-victoriametrics-operator" cannot list resource "clusterrolebindings" in API group "rbac.authorization.k8s.io" at the cluster scope
  ```
  <!-- markdownlint-enable line-length -->

* Some objects are stuck in a deleting state with a `deletionTimestamp` set and are never physically removed.

**Root cause:**

The victoriametrics-operator must handle deletion of objects carrying its `finalizer` and remove the finalizer itself. A
`Clean Install` removed the operator first, so objects deleted afterward have no controller left to process their
finalizer and stay stuck in termination.

**How to check:**

1. List the stuck objects and confirm they carry a `deletionTimestamp` and the operator finalizer:

   ```bash
   kubectl -n <namespace> get <object_type> <object_name> \
     -o jsonpath='{.metadata.deletionTimestamp}{"  "}{.metadata.finalizers}{"\n"}'
   ```

**How to fix:**

1. **DANGEROUS — removing a finalizer bypasses the operator's cleanup logic; do this only for objects already in the
   delete state with no operator left to reconcile them.** Remove the `apps.victoriametrics.com/finalizer` from each
   stuck object, or clear the whole finalizers list when there are no others:

   ```bash
   kubectl -n <namespace> patch <object_type> <object_name> \
     -p '{"metadata":{"finalizers":null}}' --type=merge
   ```

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### EFS persistence volume overflow

**Symptoms:**

* VictoriaMetrics selects data slowly or returns query results after a long delay.
* AWS CloudWatch or AWS EFS shows the EFS volume throughput fully utilized.

**Root cause:**

The EFS volume was created with the `Bursting` throughput type, which is usually not enough for VictoriaMetrics.

**How to check:**

1. In AWS CloudWatch or the EFS console, read the throughput utilization for the volume backing VictoriaMetrics.
   Sustained 100% utilization confirms the cause.

**How to fix:**

1. Switch the EFS volume to `Provisioned Throughput` (minimum 10 Mb/s, incurs cost) or at least `Elastic Throughput`.
2. Alternatively, replace AWS EFS with an AWS EBS volume of at least `Throughput Optimized HDD` type, which is cheaper
   and offers more throughput; faster SSD types also work.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### Cannot read stream body: response exceeds promscrape.maxScrapeSize

**Symptoms:**

* VMAgent logs a scrape error:

  <!-- markdownlint-disable line-length -->
  ```text
  Cannot read stream body in 1 seconds: the response from "https://x.x.x.x:6443/metrics" exceeds -promscrape.maxScrapeSize=16777216; either reduce the response size for the target or increase -promscrape.maxScrapeSize
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

VMAgent's default `promscrape.maxScrapeSize` is `16777216` bytes (16 MB). A target whose scrape response is larger is
rejected.

**How to check:**

1. Read the VMAgent log for the `exceeds -promscrape.maxScrapeSize` error and note the target URL.

**How to fix:**

1. Increase the maximum scrape size in deploy parameters and redeploy:

   ```yaml
   victoriametrics:
     vmAgent:
       extraArgs:
         promscrape.maxScrapeSize: 256MiB
   ```

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### VictoriaMetrics pods continuously restart on invalid configuration

**Symptoms:**

* VMAgent or VMAlertmanager pods restart continuously.
* The pod log names the invalid part of the configuration, for example:

  <!-- markdownlint-disable line-length -->
  ```text
  cannot read "/etc/vmagent/config_out/vmagent.env.yaml": cannot parse Prometheus config from "/etc/vmagent/config_out/vmagent.env.yaml": cannot parse `scrape_config`: cannot parse auth config for `job_name` "serviceScrape/kafka/kafka-service-monitor-jmx-exporter/0": missing `username` in `basic_auth` section
  ```
  <!-- markdownlint-enable line-length -->

  or:

  <!-- markdownlint-disable line-length -->
  ```text
  ts=2024-03-26T06:58:37.006Z caller=coordinator.go:118 level=error component=configuration msg="Loading configuration file failed" file=/etc/alertmanager/config/alertmanager.yaml err="missing to address in email config"
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

The generated configuration is invalid, and the pod refuses to start. VMAgent, VMAlert, and VMAlertmanager are the most
affected because users supply custom settings for them. The offending fragment is printed in the log. The same failure
appears for standalone Alertmanager; see also **Alertmanager will not start: "Loading configuration file failed"**. When
an `AlertmanagerConfig` fails to convert rather than the pod failing to start, see also **VictoriaMetrics config secret
does not contain inhibit rules**.

**How to check:**

1. Read the pod log and parse the error. In the first example the problem is in a `ServiceMonitor` / `VMServiceScrape`
   named `kafka-service-monitor-jmx-exporter` in the `kafka` namespace, whose `basicAuth` has no `username`. In the
   second example an email notification channel config (in `AlertmanagerConfig` / `VMAlertmanagerConfig` or the
   additional-config Secret) has no `to` address.

**How to fix:**

1. Find the CR or Secret named in the error and correct the missing field — add the `username` to the `basicAuth`
   section, or add the missing email address to the notification channel config.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### VictoriaMetrics config secret does not contain inhibit rules

**Symptoms:**

* The VictoriaMetrics configuration Secret does not contain the `inhibit_rules` you defined in an `AlertmanagerConfig`.

**Root cause:**

The victoriametrics-operator converts a Prometheus `AlertmanagerConfig` into a `VMAlertmanagerConfig`, for which `route`
and `receivers` are mandatory. They are optional for a Prometheus `AlertmanagerConfig`, so an `AlertmanagerConfig` that
omits them fails conversion and its inhibit rules never reach the VM config Secret. See also **VictoriaMetrics pods
continuously restart on invalid configuration**, the pod-startup failure when the generated config itself is invalid.

**How to check:**

1. Read the `AlertmanagerConfig` and confirm whether `route` and `receivers` are present. Their absence is the cause.

**How to fix:**

1. Add a `receivers` entry and a `route` to the `AlertmanagerConfig` so the conversion succeeds:

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
         targetMatch:
           - matchType: '=~'
             name: severity
             value: (high|warning|critical)
     receivers:
       - name: base
     route:
       receiver: base
   ```

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### VictoriaMetrics rejects writes with HTTP 503 and logs switching to read-only mode

**Symptoms:**

* vmagent (or any writer) is rejected and retries; ingestion stops.
* VMSingle logs the switch to read-only mode, and clients receive an HTTP 503 with a matching body:

  <!-- markdownlint-disable line-length -->
  ```text
  warn storage.go:669 switching the storage at /vmstorage-data to read-only mode, since it has less than -storage.minFreeDiskSpaceBytes=838860800 of free space: 838803456 bytes left
  ```
  <!-- markdownlint-enable line-length -->

  ```text
  cannot store metrics: the storage is in read-only mode; check -storage.minFreeDiskSpaceBytes command-line flag value
  ```

**Root cause:**

Free space at `-storageDataPath` dropped below `-storage.minFreeDiskSpaceBytes`. VictoriaMetrics deliberately switches
to read-only to protect the existing data and stops accepting new samples, returning HTTP 503 to writers. vmagent treats
the 503 as retriable and buffers, so no data is lost immediately, but ingestion is halted until space is freed. The
first quoted line above is from the VictoriaMetrics cluster (`vmstorage`); on the VMSingle this deployment runs the
storage path is the configured `-storageDataPath`, not `/vmstorage-data`. Match on `switching the storage at ... to
read-only mode` and the second line rather than the exact path, and treat the byte values as configuration-dependent,
not fixed literals.

**How to check:**

1. Confirm read-only state without changing anything — scrape the VMSingle `/metrics` (default port 8428) and read
   `vm_storage_is_read_only`. A value of `1` means read-only, `0` means normal.
2. Check the data volume usage:

   ```bash
   kubectl exec <vmsingle-pod> -n <namespace> -- df -h <storageDataPath>
   ```

**How to fix:**

1. Expand the PVC (EBS/EFS online expansion) so free space rises above the threshold; VictoriaMetrics leaves read-only
   automatically once space is available.
2. **DANGEROUS — lowering the retention period permanently deletes every metric older than the new period.** If
   expansion is not possible immediately, lower the retention by setting `spec.victoriametrics.vmSingle.retentionPeriod`
   on the `PlatformMonitoring` CR, so background merges drop old month-partitions and free space. Set it on the parent
   CR, not on the generated VMSingle CR: the operator overwrites the child `VMSingle.Spec` on every reconcile
   (`controllers/victoriametrics/vmsingle/handlers.go:124`), so an edit to the generated CR is reverted. Space is not
   reclaimed immediately: it frees only as whole month-partitions age out past the new retention, so the read-only
   state can persist for a while after the change. See also **Disk usage does not drop right after lowering the
   VictoriaMetrics retention period**.

**How to avoid this issue:**

Size the volume for at least 20% free space over steady-state, and alert on `vm_storage_is_read_only == 1` and on the
free-bytes trend so you act before writes stop.

**Data to collect:**

* `vm_storage_is_read_only` value and `df -h` of the storage path.
* vmagent client logs showing the 503 retry loop.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [vmstorage: excessive logging when switching to read-only mode (Issue #5159)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/5159)
* [vmagent doesn't recover (Issue #3032)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/3032)
* [VictoriaMetrics: Cluster version (read-only mode / minFreeDiskSpaceBytes)](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/)
<!-- markdownlint-enable line-length -->

### Grafana panels backed by VictoriaMetrics error with "exceeds ...; increase -search.maxUniqueTimeseries"

**Symptoms:**

* Broad queries fail while narrow ones work; Grafana shows an error bubble.
* The response body carries a limit error, for example:

  <!-- markdownlint-disable line-length -->
  ```text
  the number of matching unique timeseries exceeds 300000; either narrow down the search or increase -search.maxUniqueTimeseries
  ```
  <!-- markdownlint-enable line-length -->

* On newer builds the wording is `the number of matching timeseries exceeds N; either narrow down the search or
  increase -search.max* command-line flag values (the most likely limit is -search.maxUniqueTimeseries)`.

**Root cause:**

The query selects more unique series than `-search.maxUniqueTimeseries` allows (or exceeds `-search.maxSamplesPerQuery`
or `-search.maxQueryDuration`). These are protective caps against unexpectedly heavy queries; a dashboard with a wide
matcher or a long range trips them. The integers shown are the configured or default limits (default
`-search.maxQueryDuration` is `30s`), not fixed literals.

**How to check:**

1. Reproduce with a bounded matcher to confirm it is the series count, not connectivity — run the same query in vmui
   with a tighter label filter or shorter range and confirm it succeeds.
2. Read the exact limit named in the error string to know which flag to raise.

**How to fix:**

1. Fix the query first: add label filters, shorten the time range, or use recording rules to pre-aggregate the heavy
   expression.
2. If the limit is genuinely too low for the workload, raise the named flag (`-search.maxUniqueTimeseries`,
   `-search.maxSamplesPerQuery`, or `-search.maxQueryDuration`) through `spec.victoriametrics.vmSingle.extraArgs` on the
   `PlatformMonitoring` CR — a map of flag name (without the leading dash) to value. Set it on the parent CR, not on the
   generated VMSingle CR: the operator overwrites the child `VMSingle.Spec` on every reconcile
   (`controllers/victoriametrics/vmsingle/handlers.go:124`), so an edit to the generated CR is reverted. Higher limits
   raise the RAM/CPU cost per query.

**How to avoid this issue:**

Design dashboards against bounded label sets, and pre-aggregate with recording rules rather than querying raw
high-cardinality series over long ranges.

**Data to collect:**

* The full verbatim error string (it names the exact limit).
* The offending PromQL/MetricsQL and time range.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [How to increase search.maxUniqueTimeseries (Issue #597)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/597)
* [Error: cannot process MetricBlock / maxSamplesPerQuery (Issue #7913)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/7913)
* [VictoriaMetrics: Cluster version (-search.maxQueryDuration and limits)](https://docs.victoriametrics.com/victoriametrics/cluster-victoriametrics/)
<!-- markdownlint-enable line-length -->

### VictoriaMetrics single is OOMKilled or slow under high cardinality

**Symptoms:**

* The VMSingle pod is OOMKilled and restarts, or ingestion/query latency spikes.
* A high percentage of "slow inserts" and a rising active-series count.

**Root cause:**

A label with many unique values (or high churn from frequently changing label values) inflates the number of active
time series. VictoriaMetrics needs RAM proportional to active series; when free RAM drops too low, cache evictions cause
excessive I/O and slowdowns, and the process can be OOMKilled. The VictoriaMetrics capacity-planning docs recommend
leaving about 50% of free RAM to reduce the probability of OOM crashes and cache evictions, and about 50% spare CPU to
absorb workload spikes. This deployment runs one time-series backend at a time; for the Prometheus backend see also
**Prometheus is OOMKilled and memory climbs continuously**.

**How to check:**

1. Use the built-in cardinality explorer (vmui, Cardinality) or the official VictoriaMetrics Grafana dashboards to find
   the label/metric with the largest share of series and the churn-rate graph.
2. Check that memory is the constraint — compare pod memory usage against its limit:

   ```bash
   kubectl top pod <vmsingle-pod> -n <namespace>
   ```

**How to fix:**

1. Eliminate the high-cardinality label at ingestion with vmagent relabeling (`action: labeldrop` or replace with a
   bounded value), or use stream aggregation to pre-aggregate before storage.
2. Give the pod more RAM as short-term relief, keeping about 50% free RAM per the docs above.

**How to avoid this issue:**

Keep about 50% spare RAM and 50% spare CPU headroom, and monitor active series and churn rate on the official dashboards
so a new label is caught before it destabilizes the instance.

**Data to collect:**

* Cardinality explorer output naming the top label/metric.
* Active-series and churn-rate graphs; pod memory limit vs. usage.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [VictoriaMetrics: FAQ (high cardinality / churn rate)](https://docs.victoriametrics.com/victoriametrics/faq/)
* [VictoriaMetrics: Single-node version (spare resources)](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/)
<!-- markdownlint-enable line-length -->

### Disk usage does not drop right after lowering the VictoriaMetrics retention period

**Symptoms:**

* `-retentionPeriod` was reduced but the storage volume does not shrink for a while.
* Old data still returns in queries past the point the operator expected it gone.

**Root cause:**

VictoriaMetrics partitions data by month and deletes in bulk. A part (immutable file) is only dropped when all of its
data points fall outside the retention window, so freeing space lags the retention change by up to a partition boundary
rather than happening the instant a sample crosses the threshold.

**How to check:**

1. Confirm the configured retention took effect — check the VMSingle flags or `/metrics` for the running
   `-retentionPeriod` value.
2. Watch storage usage trend over the following hours or days to confirm space is released as partitions age out:

   ```bash
   kubectl exec <vmsingle-pod> -n <namespace> -- df -h <storageDataPath>
   ```

**How to fix:**

1. Wait for background merges to drop aged partitions — this is expected behavior, not a fault.
2. If space must be reclaimed immediately, expand the PVC as the safe interim measure rather than deleting data by hand.

**How to avoid this issue:**

Plan retention changes knowing space frees on partition boundaries, and size the volume for the transition window.

**Data to collect:**

* Running `-retentionPeriod` value and `df -h` trend over time.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Data Retention and Deduplication - VictoriaMetrics](https://wordsus.com/en/victoriametrics/data-retention-and-deduplication)
* [VictoriaMetrics: Single-node version (retention)](https://docs.victoriametrics.com/victoriametrics/single-server-victoriametrics/)
<!-- markdownlint-enable line-length -->

### vmagent's buffer volume fills up and it starts dropping the oldest samples

**Symptoms:**

* The vmagent PVC approaches full; gaps appear in data at the remote storage.
* Metrics such as `vmagent_remotewrite_pending_data_bytes` grow constantly.

**Root cause:**

When the remote storage (VMSingle) cannot keep up or is unreachable, vmagent buffers unsent data on disk at
`-remoteWrite.tmpDataPath`. The on-disk queue is capped by `-remoteWrite.maxDiskUsagePerURL`; once the cap (or the
volume) is reached, vmagent drops the oldest blocks to make room, which shows up as gaps. Buffered data is stored in
chunks of roughly 500 MB, so caps below 500 MB behave unexpectedly.

**How to check:**

1. Confirm the queue is growing and the destination is behind — query `vmagent_remotewrite_pending_data_bytes` and
   `vmagent_remotewrite_block_size_rows_dropped_total`.
2. Check the buffer volume usage:

   ```bash
   kubectl exec <vmagent-pod> -n <namespace> -- df -h <tmpDataPath>
   ```

**How to fix:**

1. Fix the downstream first — restore VMSingle capacity/availability (see the VMSingle read-only case) so vmagent can
   drain its queue.
2. Increase throughput and headroom: raise `-remoteWrite.queues` (the operator sizes `-remoteWrite.maxDiskUsagePerURL`
   from the configured StatefulStorage) and enlarge the vmagent PVC so bursts do not overflow.

**How to avoid this issue:**

Run vmagent in StatefulMode with a persistent queue volume sized for your worst expected outage, and alert on
`vmagent_remotewrite_pending_data_bytes` trending up.

**Data to collect:**

* `vmagent_remotewrite_pending_data_bytes` and dropped-rows counters over the window.
* `df -h` of the buffer volume and the configured `-remoteWrite.maxDiskUsagePerURL`.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Misleading documentation for calculating disk space for persistence queue (Issue #7055)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/7055)
* [Warn if --remoteWrite.maxDiskUsagePerURL set lower than 500MB (Issue #4195)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/4195)
* [VictoriaMetrics: vmagent (persistent queue / maxDiskUsagePerURL)](https://docs.victoriametrics.com/victoriametrics/vmagent/)
<!-- markdownlint-enable line-length -->

### vmagent returns HTTP 429 "Too Many Requests" and pushers get errors

**Symptoms:**

* Clients pushing to vmagent receive `429 Too Many Requests`.
* Data appears to stop flowing while at least one remote-write URL is unavailable.
* Only a vmagent that receives client pushes is affected; a scrape-only vmagent does not return 429 to pushers.

**Root cause:**

With `-remoteWrite.disableOnDiskQueue` set, vmagent cannot buffer to disk, so when the remote storage cannot keep up it
returns 429 to its HTTP push clients. If `-remoteWrite.dropSamplesOnOverload` is set or multiple
`-remoteWrite.disableOnDiskQueue` URLs are configured, samples are silently dropped instead of erroring. This applies
only when `-remoteWrite.disableOnDiskQueue` is set and clients push samples to vmagent over its HTTP push API; a default
vmagent that only scrapes targets, with the on-disk queue enabled, does not return 429 to pushers.

**How to check:**

1. Check the running vmagent flags for `-remoteWrite.disableOnDiskQueue` / `-remoteWrite.dropSamplesOnOverload` on
   `/metrics` and in the pod spec.
2. Confirm which remote-write URL is unhealthy by testing connectivity from the vmagent pod to each VMSingle endpoint.

**How to fix:**

1. Restore the unhealthy remote-write destination so vmagent stops shedding load.
2. If durability across outages matters more than immediate backpressure, leave the on-disk queue enabled (do not set
   `-remoteWrite.disableOnDiskQueue`) and size the queue volume accordingly.

**How to avoid this issue:**

Choose the queue policy deliberately: on-disk queue for durability, `disableOnDiskQueue` only where you prefer fast
failure and have upstream retry.

**Data to collect:**

* Running vmagent remote-write flags.
* Per-URL connectivity test results and vmagent logs at the time of the 429s.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [vmagent returns 429 if persistent queue is disabled and one RW is unavailable (Issue #9565)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/9565)
* [VictoriaMetrics: vmagent (backpressure behavior)](https://docs.victoriametrics.com/victoriametrics/vmagent/)
<!-- markdownlint-enable line-length -->

### Requests to a monitoring UI through VMAuth fail with "missing route"

**Symptoms:**

* A URL that omits a trailing slash fails; VMAuth logs and returns:

  ```text
  remoteAddr: "..."; requestURI: /alertmanager/; missing route for "/alertmanager/"
  ```

**Root cause:**

When a request path lacks a trailing slash, VMAuth issues a 301 redirect that adds the `/`, but the `src_paths` entry in
the user's `url_map` only matches the exact path without it, so the redirected `/alertmanager/` has no matching route.

**How to check:**

1. Reproduce with `curl -v` and observe the 301 to the trailing-slash path followed by the `missing route` error.
2. Inspect the VMAuth `url_map` `src_paths` for the affected user and confirm they do not match the slashed form.

**How to fix:**

1. Make the affected route a regex that matches both forms, for example `/alertmanager.*` (or `/alertmanager/.*`), in
   the user's `paths` under `spec.victoriametrics.vmUser.targetRefs` on the `PlatformMonitoring` CR — the operator
   renders those `paths` into the VMUser `src_paths`. Set them on the parent CR, not on the generated VMUser or the
   VMAuth config it produces: the operator rebuilds the VMUser `TargetRefs` from that field and overwrites the child
   `VMUser.Spec` on every reconcile (`controllers/victoriametrics/vmuser/handlers.go:28`), so a manual edit to the
   generated config is reverted. The built-in default routes carry paths hard-coded in the operator; changing one you
   cannot express through `targetRefs` requires an operator or chart change.

**How to avoid this issue:**

Author `src_paths` as regexes (`.*`) rather than exact literals so redirects and sub-paths still route.

**Data to collect:**

* `curl -v` transcript showing the 301 and `missing route`.
* The VMAuth `url_map` for the user.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [vmauth: src_paths without / results in "missing route" (Issue #4868)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/4868)
* [VictoriaMetrics: vmauth (routing / src_paths)](https://docs.victoriametrics.com/victoriametrics/vmauth/)
<!-- markdownlint-enable line-length -->

### vmalert shows alerts firing but Alertmanager never receives them

**Symptoms:**

* vmalert's UI lists a rule as firing, but no notification is produced.
* vmalert logs an error posting to the notifier, for example `Post ... EOF` or an empty reply from the notifier URL.

**Root cause:**

vmalert cannot deliver to the `-notifier.url` — the Alertmanager address/port is wrong or unreachable. An `EOF` or
`Empty reply from server` indicates it connected to something that is not the Alertmanager API, for example the wrong
port.

**How to check:**

1. Read vmalert logs for the notifier post error:

   ```bash
   kubectl logs <vmalert-pod> -n <namespace> | grep -i notifier
   ```

2. Test the notifier endpoint from the vmalert pod:

   ```bash
   kubectl exec <vmalert-pod> -n <namespace> -- wget -qO- http://<alertmanager>:9093/api/v2/status
   ```

**How to fix:**

1. Point the notifier at the correct Alertmanager service and port (`9093` for the v2 API) by setting
   `spec.victoriametrics.vmAlert.notifier` (a single notifier) or `spec.victoriametrics.vmAlert.notifiers` (a list) on
   the `PlatformMonitoring` CR, and restore connectivity if a NetworkPolicy blocks it. Set it on the parent CR, not on
   the generated VMAlert CR: the operator overwrites the child `VMAlert.Spec` on every reconcile
   (`controllers/victoriametrics/vmalert/handlers.go:124`), so an edit to the generated CR is reverted.

**How to avoid this issue:**

Verify the notifier endpoint responds on the v2 API from vmalert's network context as part of deployment validation.

**Data to collect:**

* vmalert notifier error logs.
* Connectivity test from vmalert to the Alertmanager API.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [vmalert failed to send alert to alertmanager: Post EOF (Issue #6758)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/6758)
* [VictoriaMetrics: vmalert (notifier.url)](https://docs.victoriametrics.com/victoriametrics/vmalert/)
<!-- markdownlint-enable line-length -->

### vmalert recording-rule results are not written back: "queue is full"

**Symptoms:**

* Recording-rule metrics are missing from storage.
* vmalert logs a remote-write queue error:

  <!-- markdownlint-disable line-length -->
  ```text
  group "...": rule "...": remote write failure: failed to push timeseries - queue is full (100000 entries). Queue size is controlled by -remoteWrite.maxQueueSize flag
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

vmalert's remote-write queue to `-remoteWrite.url` filled because the destination cannot keep up or the rules generate
more series than the queue can drain, so results are dropped instead of persisted.

**How to check:**

1. Read vmalert logs for the `queue is full` message and note the configured `-remoteWrite.maxQueueSize`.
2. Confirm the remote-write destination (VMSingle, or vmagent as fan-out) is healthy and keeping up.

**How to fix:**

1. Restore or scale the remote-write destination so the queue drains.
2. Increase `-remoteWrite.maxBatchSize` and, where appropriate, route through vmagent as a fan-out proxy that buffers to
   disk during storage outages.

**How to avoid this issue:**

Persist rule results through vmagent (durable on-disk queue) rather than writing directly to storage, so a storage blip
does not drop recording-rule output.

**Data to collect:**

* vmalert `queue is full` log lines and the remote-write flags.
* Health of the remote-write destination.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [vmalert ignores -remoteWrite.maxQueueSize parameter (Issue #4845)](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/4845)
* [VictoriaMetrics: vmalert (remoteWrite / fan-out via vmagent)](https://docs.victoriametrics.com/victoriametrics/vmalert/)
<!-- markdownlint-enable line-length -->

## Prometheus

### Prometheus fails to write with permission denied on hostPath volume

**Symptoms:**

* Prometheus cannot write to its data volume when using a `hostPath` PersistentVolume with a custom `securityContext`:

  ```text
  mkdir /prometheus/wal: permission denied
  ```

**Root cause:**

The Prometheus image creates the data directory as `root`, changes ownership to `nobody` (uid:gid `65534:65534`), and
runs as `nobody`. A custom user set through `securityContext` then cannot write to a host directory it does not own.

**How to check:**

1. Read the Prometheus pod log for the `permission denied` error on `/prometheus/wal`.
2. Compare the `securityContext` UID/GID against the ownership of the host directory backing the PV.

**How to fix:**

1. On the node hosting the volume, create the `prometheus-db` subdirectory and set ownership to the custom user you run
   as (example uses `2001:2001`):

   ```bash
   mkdir -p /mnt/data/prometheus/prometheus-db
   chown -R 2001:2001 /mnt/data/prometheus
   ```

2. Alternatively, run Prometheus as `nobody` (uid/gid `65534`) to match the image default.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/installation/storage.md`
<!-- markdownlint-enable line-length -->

### Prometheus is not installed even though its Install flag is true

**Symptoms:**

* Both the Prometheus stack and the VictoriaMetrics stack have `install: true`, but no Prometheus is deployed.

**Root cause:**

`FillEmptyWithDefaults` force-disables Prometheus when both Prometheus and the VictoriaMetrics operator are set to
install. The time-series backend is either VictoriaMetrics or Prometheus, and VictoriaMetrics wins: the operator sets
`Prometheus.Install` to `false` at the start of every reconcile.

**How to check:**

1. Read both install flags in the CR:

   ```bash
   kubectl -n <namespace> get platformmonitoring <cr-name> \
     -o jsonpath='prometheus={.spec.prometheus.install} vmOperator={.spec.victoriametrics.vmOperator.install}{"\n"}'
   ```

   Both `true` means Prometheus is intentionally suppressed.

**How to fix:**

1. Choose one backend. To run Prometheus, set the VictoriaMetrics operator install to `false`; to run VictoriaMetrics,
   leave Prometheus disabled. Redeploy.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `api/v1/platformmonitoring_types.go`
<!-- markdownlint-enable line-length -->

### Prometheus pod is in CrashLoopBackOff with "lock DB directory: resource temporarily unavailable"

**Symptoms:**

* The Prometheus pod restarts repeatedly (CrashLoopBackOff).
* Logs show a lock warning followed by a fatal storage error:

  <!-- markdownlint-disable line-length -->
  ```text
  caller=dir_locker.go:77 level=warn component=tsdb msg="A lockfile from a previous execution already existed. It was replaced" file=/prometheus/lock
  caller=main.go:1081 level=error err="opening storage failed: lock DB directory: resource temporarily unavailable"
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

Two live Prometheus processes hold the same TSDB data directory at once. On this deployment Prometheus runs as a
StatefulSet managed by prometheus-operator, so a routine rollout does not cause this: the StatefulSet does not start a
replacement pod of the same name on the same PVC until the old pod is deleted, and the warned-about stale lockfile from
a previous execution is automatically replaced rather than blocking startup. The fatal `lock DB directory: resource
temporarily unavailable` means another process is actively holding the lock — a second Prometheus writing to the same
TSDB, or the same data volume shared or duplicated across two pods. This is uncommon here, because a `ReadWriteOnce` PVC
bound to a StatefulSet admits only one writer; it shows up mainly when a volume is cloned, mounted read-write in two
places, or a second Prometheus is pointed at the same path.

**How to check:**

1. List the pods on the volume and confirm whether two Prometheus pods are Running against one PVC. A healthy result is
   exactly one Running pod per data volume:

   ```bash
   kubectl get pods -n <namespace> -l app.kubernetes.io/name=prometheus -o wide
   ```

2. Confirm the access mode of the PVC. `ReadWriteOnce` bound to two pods on different nodes cannot work:

   ```bash
   kubectl get pvc <claim> -n <namespace> -o jsonpath='{.spec.accessModes}'
   ```

**How to fix:**

1. Ensure only one live Prometheus process owns the data volume. Identify every pod or process mounting the TSDB path
   and stop the duplicate, so a single writer holds the lock.
2. Confirm the PVC is `ReadWriteOnce` and bound to a single pod, and that the volume is not cloned or mounted read-write
   elsewhere. If a pod is stuck terminating on the same PVC, resolve that first.

**How to avoid this issue:**

Keep Prometheus a StatefulSet with one PVC per replica and never share one TSDB volume between pods. Note that
`--storage.tsdb.no-lockfile` removes the lock guard but does not fix concurrent writers and can corrupt data, so do not
use it to work around this.

**Data to collect:**

* Full pod logs from the crashing container (include `--previous`).
* `kubectl get pods -o wide` and the PVC/PV binding for the data volume.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Prometheus: Opening storage failed: lock DB directory (Red Hat)](https://access.redhat.com/solutions/6976141)
* [CrashLoopBackoff: Disk ran out of storage? Two conflicting instances (Issue #8140)](https://github.com/prometheus/prometheus/issues/8140)
<!-- markdownlint-enable line-length -->

### Prometheus will not start after a crash: "repair corrupted WAL"

**Symptoms:**

* Prometheus crash-loops after an unclean shutdown, node crash, or out-of-space event.
* Logs show a WAL repair error, sometimes referencing a checkpoint or a padded page:

  <!-- markdownlint-disable line-length -->
  ```text
  caller=main.go:740 err="opening storage failed: repair corrupted WAL: cannot handle error: open WAL segment: 0: open /prometheus/wal/00000000: no such file or directory"
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

The write-ahead log or its checkpoint was left in an inconsistent state (disk filled, kernel OOM, or hard power loss
mid-write). Prometheus attempts automatic WAL repair on startup by truncating from the corruption point; when the
checkpoint itself is corrupt or a segment is missing, automatic repair fails and startup aborts.

**How to check:**

1. Read the startup logs and note the exact segment or checkpoint named. A healthy start logs `msg="WAL segment loaded"`
   / `msg="TSDB started"` with no `repair corrupted WAL`:

   ```bash
   kubectl logs <prometheus-pod> -n <namespace> -c prometheus --previous
   ```

2. Inspect the WAL directory size and segment list from a debug shell or ephemeral pod that mounts the volume
   (`ls -l /prometheus/wal`, `du -sh /prometheus/wal`), and confirm whether `prometheus_tsdb_wal_corruptions_total` had
   been rising before the crash.

**How to fix:**

1. First try a normal restart so Prometheus can run its own automatic truncation-based repair; recent blocks already
   compacted to disk are preserved.
2. **DANGEROUS — deleting WAL segments discards every not-yet-compacted sample (up to about 2 hours of the most recent
   metrics).** If automatic repair cannot complete, mount the volume with a temporary maintenance pod (not the crashing
   Prometheus), back up the directory, then remove only the corrupted segment named in the log and restart Prometheus:

   ```bash
   cp -a /prometheus/wal /prometheus/wal.bak && rm /prometheus/wal/<corrupted-segment>
   ```

**How to avoid this issue:**

Keep free headroom on the TSDB volume so a full disk cannot corrupt the WAL, enable WAL compression, and use faster
storage so replay/compaction finish before the next restart.

**Data to collect:**

* Startup logs naming the corrupted segment/checkpoint.
* `ls -l /prometheus/wal` and the value of `prometheus_tsdb_wal_corruptions_total` if scraped.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Crash during startup corrupted WAL (Issue #4603)](https://github.com/prometheus/prometheus/issues/4603)
* [Failing with "opening storage failed: repair corrupted WAL" (Issue #6898)](https://github.com/prometheus/prometheus/issues/6898)
<!-- markdownlint-enable line-length -->

### Prometheus restarts every few minutes and never becomes Ready after a restart

**Symptoms:**

* After a node drain, eviction, or upgrade, the Prometheus pod is Running but never Ready, then is killed and restarts.
* Logs stop at WAL replay and then show a SIGTERM before replay finishes:

  ```text
  caller=head.go:612 level=info component=tsdb msg="Replaying WAL, this may take a while"
  caller=main.go:828 level=warn msg="Received SIGTERM, exiting gracefully..."
  ```

**Root cause:**

WAL replay on startup takes longer than the liveness/readiness probe allows, so Kubernetes kills the pod before replay
completes. The next start has an even larger WAL, so the loop is self-perpetuating. High cardinality, high retention,
and slow storage make replay run for many minutes to hours.

**How to check:**

1. Count WAL segments to estimate replay time (roughly seconds per segment). A few hundred segments is normal; thousands
   means a long replay:

   ```bash
   kubectl exec <prometheus-pod> -n <namespace> -c prometheus -- sh -c 'ls -1 /prometheus/wal | wc -l'
   ```

2. Confirm the pod is being killed on the probe rather than crashing on data — run
   `kubectl describe pod <prometheus-pod> -n <namespace>` and look for `Liveness probe failed` / `Killing`.

**How to fix:**

1. Shorten the replay. Bound WAL growth by cutting cardinality and churn (see
   **Prometheus is OOMKilled and memory climbs continuously**), give the pod more memory so replay is not also fighting
   OOM, and move the TSDB to faster storage if replay is I/O-bound. This deployment does not expose
   `spec.maximumStartupDurationSeconds` — `api/v1.Prometheus` has no such field and the operator overwrites the child
   `Prometheus.Spec` on every reconcile, so editing it on the generated `Prometheus` CR is reverted.
2. If you must raise the prometheus-operator startup allowance directly, pause the operator's reconciliation of
   Prometheus first, then set `spec.maximumStartupDurationSeconds` on the generated `Prometheus` CR — but expect it to
   revert once reconciliation resumes, so treat it as a stopgap while step 1 takes effect.

**How to avoid this issue:**

Bound WAL growth by controlling cardinality and churn, and set the startup allowance above steady-state replay time so a
routine restart does not turn into a blind spot.

**Data to collect:**

* `kubectl describe pod` showing the probe failure and restart reason.
* WAL segment count and pod memory limit/usage.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [WAL Recovery Restart Loop (Issue #3391)](https://github.com/prometheus-operator/prometheus-operator/issues/3391)
* [Replaying WAL does not complete due to a restart (Issue #4585)](https://github.com/prometheus-community/helm-charts/issues/4585)
<!-- markdownlint-enable line-length -->

### Prometheus is OOMKilled and memory climbs continuously

**Symptoms:**

* The Prometheus container is OOMKilled and restarts; memory climbs steadily rather than plateauing.
* Dashboards go blank and alert evaluation lags during the memory pressure.

**Root cause:**

High cardinality or high churn. Every active series lives in the in-memory head block (roughly 1 to 4 KB per series), so
an unbounded label (pod name, user ID, request path) or heavy pod churn creates series faster than they expire. When pod
churn creates new series faster than old ones expire, the head block never shrinks and memory grows monotonically until
the process exceeds its limit. This deployment runs one time-series backend at a time; for the VictoriaMetrics backend
see also **VictoriaMetrics single is OOMKilled or slow under high cardinality**.

**How to check:**

1. Confirm head-series growth is the driver — query `prometheus_tsdb_head_series` and
   `rate(prometheus_tsdb_head_series_created_total[5m])`. A healthy instance plateaus; a cardinality problem climbs
   continuously.
2. Identify the offending metric/label with `promtool tsdb analyze /prometheus` or the `/status` TSDB stats page. Look
   for a single metric dominating the series count.

**How to fix:**

1. Drop or rewrite the offending label at scrape time with `metric_relabel_configs` (`labeldrop`, or replace the
   unbounded value with a bounded bucket) on the relevant ServiceMonitor/PodMonitor.
2. Cap each scrape with `sample_limit` / `label_limit` so a future explosion fails one scrape instead of the whole
   instance, and raise the memory limit as short-term relief.

**How to avoid this issue:**

Treat labels as a bounded API: no user IDs, emails, raw paths, or container IDs as label values. Alert on
`prometheus_tsdb_head_series` crossing a baseline threshold so growth is caught before OOM.

**Data to collect:**

* `promtool tsdb analyze` output naming top series.
* Graphs of `prometheus_tsdb_head_series` and container memory over the incident window.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [High Cardinality in Prometheus: How to Find and Fix It (Last9)](https://last9.io/blog/how-to-manage-high-cardinality-metrics-in-prometheus/)
* [How Cloudflare runs Prometheus at scale](https://blog.cloudflare.com/how-cloudflare-runs-prometheus-at-scale/)
* [Prometheus Scalability: High Cardinality and How to Fix It](https://alexandre-vazquez.com/prometheus-scalability/)
<!-- markdownlint-enable line-length -->

### Prometheus stops ingesting and the TSDB volume is full

**Symptoms:**

* Ingestion stops and gaps appear; the pod may crash-loop.
* Logs reference no space or a read-only filesystem, for example:

  ```text
  level=error err="opening storage failed: mkdir data/: read-only file system"
  ```

**Root cause:**

The TSDB PersistentVolume filled up. Prometheus needs free space for the head block, WAL, and compaction; when the
volume is full, scrape commits fail and, on some backends, the filesystem flips read-only, which then prevents a clean
restart. Retention alone does not save you if series/churn outgrew the provisioned size.

**How to check:**

1. Check volume usage from inside the pod (or a maintenance pod on the volume). Anything near 100% used is the cause:

   ```bash
   kubectl exec <prometheus-pod> -n <namespace> -c prometheus -- df -h /prometheus
   ```

2. Compare configured `--storage.tsdb.retention.time` / `.size` against the actual on-disk size (`du -sh /prometheus`).

**How to fix:**

1. Expand the PVC where the StorageClass allows it (`kubectl edit pvc <claim>` and raise
   `spec.resources.requests.storage`); AWS EBS/EFS support online expansion. This is the non-destructive fix.
2. **DANGEROUS — deleting blocks destroys the metrics stored in them permanently.** Only if expansion is impossible and
   the instance must come back, mount the volume with a maintenance pod, back up, and delete the oldest block
   directories under `/prometheus` (not the `wal/`), then restart.

**How to avoid this issue:**

Set `--storage.tsdb.retention.size` below the volume size as a hard cap, and alert on `predict_linear` of free bytes so
you expand before the volume fills.

**Data to collect:**

* `df -h /prometheus` and `du -sh /prometheus/*`.
* Retention settings and PVC/StorageClass definition.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [2.6.0: opening storage failed: read-only file system (Issue #5043)](https://github.com/prometheus/prometheus/issues/5043)
* [Oncall Adventures - When Your Prometheus Server's Disk Is Full](https://tratnayake.dev/oncall-adventures-prometheus-filled-disk)
<!-- markdownlint-enable line-length -->

## Alertmanager

### Alerts fire in Prometheus but no notification is delivered

**Symptoms:**

* An alert is visible as firing in Prometheus/vmalert but no email/Slack/webhook arrives.
* `alertmanager_notifications_failed_total` is nonzero and rising.

**Root cause:**

The receiver integration is failing — commonly a Slack webhook returning `404` (channel/URL wrong or revoked), an SMTP
auth failure, the wrong SMTP port, or a firewall blocking outbound traffic. The alert reaches Alertmanager but the send
step errors. For Gmail specifically, use port 587 with STARTTLS or port 465 with implicit SSL/TLS and authenticate with
a 16-character app password (or OAuth 2.0) rather than the account password — port 25 is not available on
`smtp.gmail.com`.

**How to check:**

1. Identify the failing integration with `rate(alertmanager_notifications_failed_total[10m])` broken down by
   `integration`.
2. Read the Alertmanager logs for the send error without changing anything:

   ```bash
   kubectl logs <alertmanager-pod> -n <namespace> | grep -i "notify\|webhook\|send"
   ```

3. Test the receiver out-of-band, for example a manual `curl -X POST` to the Slack webhook URL. A `404` back means the
   channel or URL is wrong.

**How to fix:**

1. Correct the receiver config in the AlertmanagerConfig / secret the operator renders (fix the webhook URL, use an SMTP
   app password, set port 587 for STARTTLS or 465 for SMTPS).
2. Reload/apply and re-fire a test alert to confirm delivery.

**How to avoid this issue:**

Add a meta-alert on `rate(alertmanager_notifications_failed_total[1m]) > 0` routed to an independent channel so a broken
integration is itself alerted.

**Data to collect:**

* `alertmanager_notifications_failed_total` by integration.
* Alertmanager logs around the failed send and the manual receiver test output.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Alertmanager Notification Failures (GitLab runbook)](https://runbooks.gitlab.com/monitoring/alertmanager-notification-failures/)
* [Prometheus Alertmanager Setup Guide (SMTP/Slack)](https://computingforgeeks.com/setup-prometheus-alertmanager-linux/)
<!-- markdownlint-enable line-length -->

### Alertmanager sends a literal "${SLACK_WEBHOOK}" instead of the secret value

**Symptoms:**

* Notifications fail and the logs/config show a placeholder string rather than the real webhook or password.
* The receiver rejects the request because the URL/credential is literally `${VAR}`.

**Root cause:**

Alertmanager does not expand environment variables in its configuration file. A config using
`api_url: '${SLACK_WEBHOOK}'` is sent verbatim as that string.

**How to check:**

1. Inspect the rendered config for unexpanded `${...}` tokens:

   <!-- markdownlint-disable line-length -->
   ```bash
   kubectl get secret <alertmanager-config-secret> -n <namespace> -o jsonpath='{.data.alertmanager\.yaml}' | base64 -d
   ```
   <!-- markdownlint-enable line-length -->

**How to fix:**

1. Inject the secret through an `AlertmanagerConfig` CR rather than a raw config with a shell variable. Its receiver
   fields take a Kubernetes `SecretKeySelector` — for Slack, `spec.receivers[].slackConfigs[].apiURL` referencing a key
   in a Secret — so the value is read from the Secret, never expanded from `${VAR}`. This deployment does not let you
   add arbitrary volumes or a `_file` mount to the managed Alertmanager pod, so the file-based (`api_url_file`) approach
   is not available here.

**How to avoid this issue:**

Configure every receiver secret through `AlertmanagerConfig` `SecretKeySelector` fields, so nothing depends on
shell-style expansion Alertmanager does not perform.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [How to Debug Alertmanager Issues (env vars not expanded)](https://oneuptime.com/blog/post/2026-01-28-debug-alertmanager-issues/view)
* [Alertmanager configuration (prometheus.io)](https://prometheus.io/docs/alerting/latest/configuration/)
<!-- markdownlint-enable line-length -->

### Alertmanager will not start: "Loading configuration file failed"

**Symptoms:**

* Alertmanager crash-loops or refuses to reload.
* Logs show a fatal config error naming an unknown or misplaced field:

  ```text
  level=error msg="Loading configuration file failed: unknown field ..."
  ```

**Root cause:**

An invalid or mis-scoped field in `alertmanager.yml` (for example putting `http_config` at the wrong level, or a typo in
a receiver block). Alertmanager validates the whole file on load and refuses to run with an invalid config.

**How to check:**

1. Validate the config offline with `amtool check-config alertmanager.yml`, or read the exact field named in the startup
   log.

**How to fix:**

1. Correct the offending field per the configuration reference and re-apply the AlertmanagerConfig/secret; a valid file
   lets the pod start and reload.

**How to avoid this issue:**

Run `amtool check-config` in CI before applying any Alertmanager configuration change.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Alertmanager Configuration global http_config (Issue #1274)](https://github.com/prometheus/alertmanager/issues/1274)
* [Alertmanager configuration (prometheus.io)](https://prometheus.io/docs/alerting/latest/configuration/)
<!-- markdownlint-enable line-length -->

## Grafana

### GrafanaDashboard rejected because it exceeds the etcd 1 MB limit

**Symptoms:**

* A `GrafanaDashboard` CR (or a ConfigMap holding a dashboard) is rejected for exceeding the per-resource size limit.

**Root cause:**

etcd limits each stored resource to 1 MB. A typical dashboard is 100-300 KB, but `kubectl apply` also writes the whole
object into the `kubectl.kubernetes.io/last-applied-configuration` annotation, doubling the stored size to 200-600 KB or
more and pushing large dashboards over the limit. See also **PlatformMonitoring or component CRD apply fails on
annotation or etcd size limit**, the same annotation-and-size mechanism for the operator's own CRDs.

**How to check:**

1. Measure the dashboard JSON size and check whether the object carries the `last-applied-configuration` annotation:

   ```bash
   kubectl -n <namespace> get grafanadashboard <name> \
     -o jsonpath='{.metadata.annotations.kubectl\.kubernetes\.io/last-applied-configuration}{"\n"}' | wc -c
   ```

**How to fix:**

1. Reduce the dashboard size, or host the dashboard JSON on external storage (for example Nexus) and reference it by URL
   in the `GrafanaDashboard` so the CR itself stays small:

   ```yaml
   apiVersion: grafana.integreatly.org/v1beta1
   kind: GrafanaDashboard
   metadata:
     name: helm-example-dashboard-by-url
     labels:
       app.kubernetes.io/component: monitoring  # Mandatory label
   spec:
     instanceSelector:
       matchLabels:
         dashboards: grafana
     url: <dashboard-json-url>
   ```

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/configuration.md`
<!-- markdownlint-enable line-length -->

### Custom Grafana plugins never reach Grafana

**Symptoms:**

* Custom plugins installed through the grafana-plugins-init container do not show up in Grafana; panels report a plugin
  is not found.

**Root cause:**

The `grafana-plugins-init` container installs the bundled plugins into the shared `grafana-plugins` volume it mounts at
`/opt/plugins`, and the operator mounts that same volume into the Grafana container at `/var/lib/grafana/plugins`
(`controllers/grafana/manifest.go`). The operator creates the volume and both mounts whenever the plugins init
container is configured, and it restores the Grafana spec on every reconcile, so a missing mount is not the reproducible
failure. The copy fails instead when the init container (uid `65534`) cannot write to `/opt/plugins`, or Grafana loads
nothing because the plugin is incompatible with the running Grafana version — for example a plugin built for an older
plugin API against Grafana 11.x.

**How to check:**

1. Read the `grafana-plugins-init` container log. Success prints `Plugins are successfully copied`; a
   `cp: permission denied` on `/opt/plugins` points at volume permissions, not a missing mount.
2. If the copy succeeded but a panel still reports a plugin is not found, read the Grafana container log for a plugin
   load or signature error and confirm the plugin supports the running Grafana version (11.x).

**How to fix:**

1. If the copy fails on permissions, set the pod security context (for example `fsGroup`) so uid `65534` can write to
   the `grafana-plugins` volume, then restart the pod. Do not add the volume or mount by hand — the operator already
   creates both and reverts manual pod edits on the next reconcile.
2. If the plugin is incompatible with the running Grafana version, replace it with a build that supports Grafana 11.x,
   or pin a Grafana image the plugin supports.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-grafana-plugins-init](https://github.com/Netcracker/qubership-grafana-plugins-init) (not yet confirmed on a live install).
* Derived from source: [entrypoint.sh](https://github.com/Netcracker/qubership-grafana-plugins-init/blob/main/entrypoint.sh),
  [Dockerfile](https://github.com/Netcracker/qubership-grafana-plugins-init/blob/main/Dockerfile) (runs as uid 65534)
<!-- markdownlint-enable line-length -->

### Grafana dashboards are empty and the datasource test returns "Bad Gateway"

**Symptoms:**

* Adding or testing the Prometheus/VM datasource fails; panels show no data.
* Grafana logs a proxy error:

  ```text
  http: proxy error: dial tcp 10.0.226.40:9090: i/o timeout ... status=502
  ```

**Root cause:**

Grafana's backend cannot reach the datasource URL. For the operator-managed default datasource the URL is correct by
construction — the operator renders it from the in-cluster service name and port — so a Bad Gateway there means the
backend is down (the TSDB pod is not Ready) or a NetworkPolicy or auth layer (oauth2-proxy/VMAuth) in front of the
datasource is rejecting Grafana's request. A wrong service name or port is possible only on a user-created
(non-operator) `GrafanaDatasource`. A 504 with slow-then-failing panels is a different root cause; see also **Grafana
panels time out with "504 Gateway Time-out" against Prometheus/VM**.

**How to check:**

1. Test reachability from inside the Grafana pod. A healthy backend responds; a refused/timeout confirms the network
   path:

   ```bash
   kubectl exec <grafana-pod> -n <namespace> -- wget -qO- http://<datasource-service>:<port>/-/ready
   ```

2. Verify the datasource URL configured in Grafana matches the in-cluster service name and port.

**How to fix:**

1. For the operator-managed default datasource, do not edit the datasource CR — its URL is correct by construction, and
   the operator overwrites the child `GrafanaDatasource.Spec` on every reconcile
   (`controllers/grafana/handlers.go:111-124`), so an edit is reverted. Restore the backend instead: bring the TSDB pod
   back to Ready (VMSingle default `8428`, Prometheus `9090`).
2. If a NetworkPolicy or auth proxy blocks the path, allow Grafana→backend traffic or point Grafana at the internal
   (unauthenticated) service.
3. Only for a user-created (non-operator) `GrafanaDatasource`, correct the URL to the right in-cluster service and port
   directly on that CR.

**How to avoid this issue:**

Use stable in-cluster service DNS names in the datasource, and keep Grafana's egress path to the TSDB open in the
NetworkPolicy.

**Data to collect:**

* Grafana proxy-error log line.
* Reachability test output from the Grafana pod and the datasource URL.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [502 error adding data source (prometheus) - Grafana community](https://community.grafana.com/t/502-error-adding-data-source-prometheus/15429)
* [Troubleshoot Prometheus data source issues (Grafana docs)](https://grafana.com/docs/grafana/latest/datasources/prometheus/troubleshooting/)
<!-- markdownlint-enable line-length -->

### Grafana panels time out with "504 Gateway Time-out" against Prometheus/VM

**Symptoms:**

* Panels load slowly then error, often at a fixed ~30 s boundary.
* Users see `504 Gateway Time-out` or a query timeout message.

**Root cause:**

The query takes longer than a timeout in the path — Grafana's datasource timeout, the TSDB's own query timeout, or an
intervening proxy/ingress read timeout. Heavy queries over long ranges or high cardinality are the usual trigger. A
502/Bad Gateway with empty panels is a connectivity problem, not a timeout; see also **Grafana dashboards are empty and
the datasource test returns "Bad Gateway"**.

**How to check:**

1. Run the panel's query directly against the datasource (vmui or `promtool`/curl) and measure how long it takes; if it
   exceeds ~30 s the timeout is the query, not the network.
2. Check whether an ingress/proxy in front imposes its own read timeout.

**How to fix:**

1. Optimize the query (shorter range, label filters, recording rules) so it returns within the timeout.
2. Where the workload legitimately needs longer, raise the datasource Timeout in Grafana and the TSDB query-duration
   limit (and any proxy read timeout) together.

**How to avoid this issue:**

Back dashboards with recording rules and bounded ranges; align Grafana, TSDB, and proxy timeouts intentionally.

**Data to collect:**

* Direct query timing against the datasource.
* The timeout values configured in Grafana, the TSDB, and any proxy.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Grafana timing out when querying Prometheus datasource (community)](https://community.grafana.com/t/grafana-timing-out-when-querying-prometheus-datasource/23167)
* [Troubleshoot Prometheus data source issues (Grafana docs)](https://grafana.com/docs/grafana/latest/datasources/prometheus/troubleshooting/)
<!-- markdownlint-enable line-length -->

### A GrafanaDashboard or GrafanaDatasource never appears, status "NoMatchingInstance"

**Symptoms:**

* The dashboard/datasource is not created in Grafana despite the CR existing.
* The CR status reports no matching instance:

  <!-- markdownlint-disable line-length -->
  ```text
  "message": "None of the available Grafana instances matched the selector, skipping reconciliation", "reason": "EmptyAPIReply", "type": "NoMatchingInstance"
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

In grafana-operator v5 the resource chooses its Grafana via `instanceSelector`, matched against labels on the Grafana
CR. If the labels do not match (or the target Grafana is in another namespace without `allowCrossNamespaceImport: true`,
or the Grafana instance is not Ready), the resource is skipped. `instanceSelector` is immutable once set.

**How to check:**

1. Read the CR status for the `NoMatchingInstance` condition:

   ```bash
   kubectl get grafanadashboard <name> -n <namespace> -o jsonpath='{.status}'
   ```

2. Compare the `instanceSelector.matchLabels` against the labels actually on the Grafana CR
   (`kubectl get grafana <name> -n <namespace> --show-labels`), and confirm the Grafana instance is in a ready state.

**How to fix:**

1. Align the labels: set `spec.grafana.labels` on the `PlatformMonitoring` CR so the Grafana CR carries the labels the
   resource's `instanceSelector` matches (for example `dashboards: grafana`). Set them on the parent CR, not on the
   generated Grafana CR: the operator overwrites the child `Grafana.Spec` and labels on every reconcile
   (`controllers/grafana/handlers.go:55`), so labels added to the generated CR are reverted.
2. For cross-namespace targeting, set `allowCrossNamespaceImport: true` on the resource. If `instanceSelector` itself is
   wrong, recreate the resource (it is immutable).

**How to avoid this issue:**

Standardize one instance label (for example `dashboards: grafana`) across the Grafana CR and every dashboard/datasource
`instanceSelector` the operator renders.

**Data to collect:**

* The CR `.status` block.
* Labels on the Grafana CR and the resource's `instanceSelector`.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [GrafanaDataSources stays in status NoMatchingInstance (Issue #2280)](https://github.com/grafana/grafana-operator/issues/2280)
* [Resource Selection (grafana-operator DeepWiki)](https://deepwiki.com/grafana/grafana-operator/4.2-resource-selection)
<!-- markdownlint-enable line-length -->

### Report PNGs are blank or time out with "Error while waiting for the panels to load"

**Symptoms:**

* Rendered images for reports/alerts are blank or fail intermittently.
* The renderer logs a timeout:

  <!-- markdownlint-disable line-length -->
  ```text
  "err":"TimeoutError: Waiting failed: 60000ms exceeded ...","level":"error","message":"Error while waiting for the panels to load"
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

Either the renderer captured the page before the panel's async queries finished — a slow datasource makes the dashboard
exceed the rendering timeout — or the render token lacks datasource access, or a panel plugin fails to load in headless
Chrome, producing a blank image.

**How to check:**

1. Read the renderer logs for whether it timed out waiting for panels or failed a data request (`net::ERR_*`).
2. Reproduce with a single fast panel to isolate whether it is a specific slow or broken panel.

**How to fix:**

1. Raise both Grafana's `[rendering]` timeout and the renderer's `RENDERING_TIMEOUT` above the dashboard's real load
   time, and fix the slow queries behind the panels.
2. Ensure the renderer runs with `--no-sandbox` in-container and has the shared libraries Chromium needs so the browser
   starts at all.

**How to avoid this issue:**

Keep report dashboards lightweight (bounded ranges, recording rules) and set renderer/Grafana timeouts with headroom
over steady-state render time.

**Data to collect:**

* Renderer debug logs showing the timeout or `net::ERR_*`.
* Grafana `[rendering]` and renderer timeout settings.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Full page screenshots not working, "Error while waiting for the panels to load" (Issue #488)](https://github.com/grafana/grafana-image-renderer/issues/488)
* [Troubleshoot image rendering (Grafana docs)](https://grafana.com/docs/grafana/latest/setup-grafana/image-rendering/troubleshooting/)
<!-- markdownlint-enable line-length -->

## Authentication and access

### Login loops and "Unable to find a valid CSRF token" behind oauth2-proxy

**Symptoms:**

* After authenticating with the IdP, the browser bounces back to login repeatedly.
* The error page reads:

  ```text
  Login Failed: Unable to find a valid CSRF token. Please try again.
  ```

* Logs show `AuthFailure Invalid authentication via OAuth2: unable to obtain CSRF cookie`.

**Root cause:**

The CSRF cookie set at the start of the flow is not present at the callback. Common causes: a cross-subdomain redirect
where the cookie scope does not cover the callback host, the CSRF cookie exceeding size limits when the proxy sits
behind nginx that copies only the first `Set-Cookie`, or a callback host outside the cookie's domain. On this deployment
the oauth2-proxy config is rendered fixed — the operator's `OAuthProxy` CR exposes only `image`
(`api/v1/platformmonitoring_types.go`), and
`charts/qubership-monitoring-operator/templates/oauth2-configs/secret-oauth2-proxy-config.yaml` hard-codes the config,
including `cookie_secure = false`. The cookie and domain knobs (`cookie_domains`, `whitelist_domains`,
`cookie_samesite`, `cookie_secure`, `cookie_csrf_per_request`) are therefore not settable through supported values, and
a manual edit to the rendered Secret is overwritten on the next Helm render.

**How to check:**

1. Read the oauth2-proxy logs to see which cookie is missing at callback:

   ```bash
   kubectl logs <oauth2-proxy-pod> -n <namespace> | grep -i csrf
   ```

2. In browser devtools, watch whether `_oauth2_proxy_csrf` is set on `/oauth2/start` and still present at
   `/oauth2/callback`.

**How to fix:**

1. Fix the redirect topology so the fixed config works: put the application and oauth2-proxy on hosts that share a
   cookie scope, and make sure the callback host the IdP redirects to matches the host that set the CSRF cookie. This is
   the durable remedy here, because the cookie flags themselves are not settable on this deployment.
2. If the topology cannot be aligned and a cookie setting genuinely must change (for example `cookie_samesite` or
   `cookie_csrf_per_request`), raise a chart change to expose that setting or to correct the rendered
   `secret-oauth2-proxy-config.yaml`. Do not edit the rendered Secret in place — the next Helm render overwrites it.

**How to avoid this issue:**

Keep the application and oauth2-proxy on domains that share a cookie scope, so the fixed config's cookies survive the
redirect without per-deployment cookie tuning.

**Data to collect:**

* oauth2-proxy logs around the failed callback.
* The rendered oauth2-proxy config and the browser cookie trace.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Login Failed: Unable to find a valid CSRF token (Issue #2965)](https://github.com/oauth2-proxy/oauth2-proxy/issues/2965)
* [Unable to find a valid CSRF token when redirected from another site (Issue #2830)](https://github.com/oauth2-proxy/oauth2-proxy/issues/2830)
<!-- markdownlint-enable line-length -->

## Exporters

### Not all node-exporter endpoints have UP status in OpenShift

**Symptoms:**

* On OpenShift, some node-exporter endpoints are not `UP`; requests to the node-exporter port are rejected.

**Root cause:**

OpenShift's iptables service allows only a port range (by default `30000`-`32999`). node-exporter runs on port `9900`,
outside that range, so iptables rejects requests to it.

**How to check:**

1. Confirm node-exporter listens on `9900` and that the port is outside the OpenShift iptables allowed range.
2. From a node, confirm requests to port `9900` are refused.

**How to fix:**

1. On each virtual machine except balancer nodes, open port `9900`. Edit `/etc/sysconfig/iptables` and add this line
   before `COMMIT`:

   ```text
   -A OS_FIREWALL_ALLOW -p tcp -m state --state NEW -m tcp --dport 9900 -j ACCEPT
   ```

2. **DANGEROUS — restarting iptables briefly interrupts all firewall filtering and existing connections on the node.**
   Restart the iptables service to apply the rule:

   ```bash
   systemctl restart iptables
   ```

**How to avoid this issue:**

For OpenStack, open port `9900` in the Security Group. monitoring-operator uses the OpenShift-provided
SecurityContextConstraint for node-exporter.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### Network latency exporter reports no targets or unreachable nodes

**Symptoms:**

* The exporter runs but discovers no target nodes, or every node reports as unreachable with zeroed RTT.
* Logs contain one of:

  ```text
  Error getting cluster nodes
  Failed to create node watcher
  Failed to run mtr process
  Process timeout
  ```

**Root cause:**

There are three common causes: the exporter's ServiceAccount lacks `nodes` `list`/`watch` permission, so node discovery
finds nothing (`Error getting cluster nodes`, `Failed to create node watcher`); the ports in `checkTarget` (default
`UDP:80,TCP:80,ICMP`) are not open on the nodes, so probes report unreachable with zeroed RTT (port `1` is only the
fallback when a `checkTarget` entry omits its port); or `mtr` cannot run — the binary is missing or the process
lacks root to open a raw socket (`Failed to run mtr process`), or the probe exceeds its deadline (`Process timeout`).
Empty dashboards on otherwise-healthy exporter pods share the same root: the exporter needs root (uid 0) to open raw
sockets, and on OpenShift `v4.x` root alone is not enough — it also needs the privileged RBAC flag.

**How to check:**

1. Read the exporter log for the messages above to tell the three causes apart.
2. Confirm the ServiceAccount has `nodes` `list`/`watch`:

   ```bash
   kubectl auth can-i list nodes \
     --as=system:serviceaccount:<namespace>:<exporter-serviceaccount>
   ```

3. Confirm the ports in `checkTarget` (default `UDP:80,TCP:80,ICMP`) are open on the nodes.

**How to fix:**

1. Grant the exporter ServiceAccount `nodes` `list`/`watch` so node discovery works.
2. Open the `checkTarget` ports on each node (default `UDP:80,TCP:80,ICMP`), or set `checkTarget` to ports that are
   open.
3. Ensure the `mtr` binary is present and the exporter runs as root — set `runAsUser: 0` (a numeric integer; the
   subchart schema types this field as an integer) so it can open raw sockets. On OpenShift `v4.x`, additionally set
   `.Values.networkLatencyExporter.rbac.privileged` to `true`.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-network-latency-exporter](https://github.com/Netcracker/qubership-network-latency-exporter) (not yet confirmed on a live install).
* Derived from source: [pkg/collector/node_collector.go](https://github.com/Netcracker/qubership-network-latency-exporter/blob/main/pkg/collector/node_collector.go), [pkg/collector/discover.go](https://github.com/Netcracker/qubership-network-latency-exporter/blob/main/pkg/collector/discover.go), [cmd/main.go](https://github.com/Netcracker/qubership-network-latency-exporter/blob/main/cmd/main.go)
<!-- markdownlint-enable line-length -->

### version-exporter reports scrape errors and no version metrics for a target

**Symptoms:**

* Version metrics for a target are missing while the exporter keeps serving `/metrics`.
* The self-metric `version_scrape_last_scrape_error` is `1`, and `version_scrape_scrape_errors_total{collector=...}`
  increases.

**Root cause:**

A collector failed to scrape its target and skipped the metric rather than crashing. Common causes by collector:

1. HTTP collector: the endpoint returned a non-200 status. Only HTTP `200` is accepted — any 3xx/401/403/500 is a
   failure (`Response status code is not acceptable`). Wrong content type or a JSONPath that does not match the
   configured labels also fails.
2. ConfigMap collector: missing cluster RBAC. `Failed to get namespaces` means the ServiceAccount lacks cluster
   get/list/watch on namespaces, configmaps, and secrets.
3. Postgres/SSH collectors: bad host, port, credentials, or network path.

**How to check:**

1. Query the self-metrics to confirm a scrape error and identify the failing collector:

   ```bash
   curl -s http://<version-exporter>:<port>/metrics | grep version_scrape
   ```

2. Read the exporter pod log for the collector-specific error (for example `Response status code is not acceptable` or
   `Failed to get namespaces`).

**How to fix:**

1. HTTP collector: point it at an endpoint that returns HTTP `200` with a supported content type (`application/json` or
   `text/plain`), and correct the JSONPath so its result count matches the configured labels.
2. ConfigMap collector: grant the ServiceAccount cluster get/list/watch on namespaces, configmaps, and secrets.
3. Postgres/SSH collectors: correct the host, port, credentials, and network reachability.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-version-exporter](https://github.com/Netcracker/qubership-version-exporter) (not yet confirmed on a live install).
* Derived from source: [pkg/collector/http_collector.go](https://github.com/Netcracker/qubership-version-exporter/blob/main/pkg/collector/http_collector.go), [pkg/collector/configmap_collector.go](https://github.com/Netcracker/qubership-version-exporter/blob/main/pkg/collector/configmap_collector.go)
<!-- markdownlint-enable line-length -->

### Filesystem metrics are missing and node-exporter logs a mounts read error

**Symptoms:**

* `node_filesystem_*` metrics are absent for a node.
* Logs show the filesystem collector failing to read mount info:

  ```text
  /proc/1/mounts can not be read
  ```

**Root cause:**

The operator's node-exporter DaemonSet always bind-mounts the host root at `/host` and passes `--path.rootfs=/host`, so
the host mounts are present. The filesystem collector still fails when it runs as a non-root user that lacks permission
to read the host mount table (`/proc/1/mounts`). Newly mounted filesystems that appear after node-exporter started are
also not picked up until it re-reads mounts.

**How to check:**

1. Read the node-exporter logs for the collector error:

   ```bash
   kubectl logs <node-exporter-pod> -n <namespace> | grep -i "mounts\|filesystem"
   ```

2. Confirm the container has permission to read the host mount table. The DaemonSet already mounts the host root at
   `/host` (`--path.rootfs=/host`), so the failure is a read-permission problem, not a missing mount.

**How to fix:**

1. Run the exporter with enough privilege to read the host mount table — grant the container read access to
   `/proc/1/mounts` (for example, run as root or add the capability it needs). The operator already sets the host mounts
   and `--path.rootfs=/host`, so no mount change is required.

**How to avoid this issue:**

Run node-exporter with sufficient privilege to read the host mount table; the operator's DaemonSet already provides the
standard host mounts and `--path.rootfs=/host`.

**Data to collect:**

* node-exporter logs with the collector error.
* The DaemonSet's volumeMounts and args.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [filesystem collector failed for non-root user: /proc/1/mounts can not be read (Issue #1171)](https://github.com/prometheus/node_exporter/issues/1171)
* [Cannot proper monitor rootfs (Issue #1585)](https://github.com/prometheus/node_exporter/issues/1585)
<!-- markdownlint-enable line-length -->

### No node metrics for control-plane / tainted nodes

**Symptoms:**

* Node dashboards have data for worker nodes but nothing for control-plane nodes.
* No node-exporter pod is scheduled on the tainted nodes.

**Root cause:**

The operator's node-exporter DaemonSet ships a blanket toleration (`tolerations: [{operator: Exists}]`) by default,
which tolerates every taint, so it normally schedules on control-plane and other tainted nodes. This symptom appears
only when the deployment overrode that default through `spec.nodeExporter.tolerations` with a narrower set that does not
cover the control-plane taint (for example `node-role.kubernetes.io/control-plane:NoSchedule`). The operator replaces
the DaemonSet tolerations with whatever the CR specifies, so a narrower list leaves those nodes without a node-exporter
pod.

**How to check:**

1. Confirm the missing pods with `kubectl get pods -n <namespace> -o wide -l app=node-exporter` and compare node
   coverage against `kubectl get nodes`.
2. Read the node taints:

   ```bash
   kubectl describe node <control-plane-node> | grep Taints
   ```

3. Check whether the CR overrides the default tolerations:

   ```bash
   kubectl -n <namespace> get platformmonitoring <cr-name> \
     -o jsonpath='{.spec.nodeExporter.tolerations}{"\n"}'
   ```

   Empty output means the default blanket toleration is in force; a narrower list is the likely cause.

**How to fix:**

1. Restore the blanket toleration: remove the `spec.nodeExporter.tolerations` override so the default `operator: Exists`
   is used again, or include an `operator: Exists` toleration in your override, then redeploy so pods schedule on every
   eligible node.

**How to avoid this issue:**

Keep the default `operator: Exists` toleration (or include it in any override) so node-exporter covers every node, and
alert on tainted nodes missing a DaemonSet pod.

**Data to collect:**

* `kubectl get pods -o wide` vs. `kubectl get nodes`.
* Node taints and the DaemonSet tolerations.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [How to Configure DaemonSet Tolerations for Running on Tainted Nodes](https://oneuptime.com/blog/post/2026-02-09-daemonset-tolerations-tainted-nodes/view)
* [node_exporter DaemonSet not scheduled after taints (Issue #670)](https://github.com/prometheus/node_exporter/issues/670)
<!-- markdownlint-enable line-length -->

### kube-state-metrics is OOMKilled and restarts on a large cluster

**Symptoms:**

* The kube-state-metrics pod shows `Last State: Terminated / Reason: OOMKilled` and restarts.
* `kube_*` metrics go stale during the restarts.

**Root cause:**

Memory use scales with the number of objects (pods, secrets, namespaces) KSM watches. On large or high-object clusters a
single instance exceeds its memory limit. A common rule of thumb is about 10 MB per 1,000 pods, plus overhead for
secrets and other objects.

**How to check:**

1. Confirm the kill reason:

   ```bash
   kubectl describe pod -l app.kubernetes.io/name=kube-state-metrics -n <namespace>
   ```

2. Gauge cluster size — object counts for pods/secrets/namespaces (`kubectl get pods -A | wc -l`, and similar).

**How to fix:**

1. Raise the memory limit as immediate relief, through the `kubeStateMetrics.resources` limits on the
   `PlatformMonitoring` CR.
2. Reduce the object and series count with the scope fields the CR exposes: set `kubeStateMetrics.namespaces` to a
   comma-separated namespace list, `kubeStateMetrics.scrapeResources` to the resource kinds you actually need, and
   `kubeStateMetrics.metricLabelsAllowlist` to limit the label metrics collected. The operator does not expose
   horizontal sharding (`--shard` / `--total-shards`) or the raw `--resources` / `--namespaces` / `--metric-allowlist`
   flags, so scope through these CR fields rather than the flags.

**How to avoid this issue:**

Size memory from object count, and narrow scope with the `kubeStateMetrics.namespaces`,
`kubeStateMetrics.scrapeResources`, and `kubeStateMetrics.metricLabelsAllowlist` fields before the cluster grows past a
single instance's headroom.

**Data to collect:**

* `kubectl describe pod` showing OOMKilled and the memory limit.
* Object counts and, if enabled, KSM self-metrics.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Repeated OOM'ing due to a large number of namespaces (Issue #493)](https://github.com/kubernetes/kube-state-metrics/issues/493)
* [kube-state-metrics Troubleshooting (OOMKilled / rule of thumb)](https://kubestatemetrics.com/blog/troubleshooting/)
<!-- markdownlint-enable line-length -->

### kube-state-metrics stops exposing metrics after a Kubernetes upgrade with "is forbidden"

**Symptoms:**

* Some or all `kube_*` metrics disappear, often after a cluster upgrade or chart change.
* Logs show RBAC errors:

  <!-- markdownlint-disable line-length -->
  ```text
  error listing *v1.Pod: pods is forbidden: User "system:serviceaccount:monitoring:kube-state-metrics" cannot list resource "pods" in API group "" at the cluster scope
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

Either the ServiceAccount's ClusterRole is missing a resource/verb KSM needs — commonly after a Kubernetes version bump
introduces new API versions/objects, or when the chart's ClusterRole drifts from the running version — or, with
autosharding, the ClusterRole is missing the `get` verb on `pods` that the StatefulSet-based sharding logic requires.

**How to check:**

1. Read the KSM logs for the exact forbidden resource:

   ```bash
   kubectl logs deploy/kube-state-metrics -n <namespace> | grep forbidden
   ```

2. Inspect the bound ClusterRole with `kubectl describe clusterrole kube-state-metrics` and compare against the resource
   named in the error.

**How to fix:**

1. Update the operator or chart so the rendered KSM ClusterRole carries the missing resource/verb (`list`, `watch`, and
   for autosharding `get` on `pods`), and re-apply it. The durable fix has to come from the manifest the operator
   renders: the reconciler replaces every rule on the KSM ClusterRole from its embedded manifest on each reconcile
   (`controllers/kubestatemetrics/handlers.go:57-63`), so a rule added to the ClusterRole by hand is reverted.
2. As a temporary unblock only, add the missing resource/verb to the ClusterRole directly to restore metrics until the
   operator or chart update lands — expect the next reconcile to revert it.

**How to avoid this issue:**

Pin the KSM version in the operator's values and review the changelog/RBAC before Kubernetes upgrades.

**Data to collect:**

* KSM logs naming the forbidden resource.
* The current ClusterRole and ClusterRoleBinding.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [kube-state-metrics RBAC issue: forbidden (Issue #1908)](https://github.com/prometheus-operator/prometheus-operator/issues/1908)
* [Example clusterrole needs RBAC update for GET permissions (Issue #1613)](https://github.com/kubernetes/kube-state-metrics/issues/1613)
<!-- markdownlint-enable line-length -->

### An HTTPS probe reports "probe_success 0" with a TLS handshake failure

**Symptoms:**

* `probe_success 0` for an HTTPS target that works in a browser.
* Probe logs show:

  <!-- markdownlint-disable line-length -->
  ```text
  level=ERROR source=http.go:474 msg="Error for HTTP request" err="Get \"https://.../\": remote error: tls: handshake failure"
  ```
  <!-- markdownlint-enable line-length -->

**Root cause:**

The TLS handshake between blackbox-exporter and the target fails — commonly a TLS version/cipher mismatch (the target
requires a version the module does not offer), an SNI/hostname issue, or the module verifying a certificate the exporter
does not trust.

**How to check:**

1. Re-run the probe with debug — `GET /probe?target=<url>&module=<module>&debug=true` — and read the TLS section.
2. Compare the module's `tls_config` (min/max version, `insecure_skip_verify`) against what the target requires.

**How to fix:**

1. Adjust the module's `tls_config` to offer a compatible TLS version and set the correct SNI/CA, in the blackbox module
   config the operator renders. This is the permanent fix.
2. **DANGEROUS — disables TLS certificate verification for this probe; use only to confirm the cause, never as the
   permanent fix.** Set `insecure_skip_verify: true` to confirm the handshake is the cause, then revert it once the real
   TLS settings are corrected.

**How to avoid this issue:**

Define per-target modules with explicit TLS settings that match each endpoint, and monitor
`probe_ssl_earliest_cert_expiry` so cert problems are caught early.

**Data to collect:**

* `debug=true` probe output for the target.
* The module's `tls_config` and the target's supported TLS versions.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [remote error: tls: handshake failure on 0.26 blackbox version (Issue #1390)](https://github.com/prometheus/blackbox_exporter/issues/1390)
* [remote error: tls: handshake failure (Issue #1261)](https://github.com/prometheus/blackbox_exporter/issues/1261)
<!-- markdownlint-enable line-length -->

### ICMP probes always fail with "socket: operation not permitted"

**Symptoms:**

* `probe_success 0` for every ICMP target.
* Probe logs show a socket permission error:

  ```text
  level=error msg="Error listening to socket" err="listen ip4:icmp 0.0.0.0: socket: operation not permitted"
  ```

**Root cause:**

ICMP requires raw-socket privileges. In a container without `CAP_NET_RAW` (or the right unprivileged-ICMP sysctl), the
exporter cannot open the ICMP socket, so every ICMP probe fails.

**How to check:**

1. Read the probe logs for `operation not permitted` / `socket: permission denied` on socket creation.
2. Inspect the pod's securityContext capabilities:

   ```bash
   kubectl get pod <blackbox-pod> -n <namespace> -o yaml | grep -A5 capabilities
   ```

**How to fix:**

1. Grant `CAP_NET_RAW` to the blackbox-exporter container (add it in the pod's securityContext capabilities) so it can
   open ICMP sockets.

**How to avoid this issue:**

Provision `CAP_NET_RAW` wherever ICMP modules are used, and confirm ICMP probes succeed after any securityContext
change.

**Data to collect:**

* Probe logs with the socket error.
* The pod securityContext capabilities.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Prometheus, BlackBox exporter ICMP issue (Google Groups)](https://groups.google.com/g/prometheus-users/c/jJtvd8-E8UU)
* [New icmp handling with capabilities breaks icmp checks (Issue #2360)](https://github.com/prometheus-community/helm-charts/issues/2360)
<!-- markdownlint-enable line-length -->

### cert-exporter exposes only "cert_exporter_error_total" and no certificate metrics

**Symptoms:**

* The exporter's `/metrics` endpoint is UP in Prometheus but the only series present is:

  ```text
  cert_exporter_error_total
  ```

* No `cert_exporter_*_expires_in_seconds` metrics appear.

**Root cause:**

The exporter found no certificates to parse — its include glob / annotation selector did not match any files or secrets,
or it lacks permission to read the referenced certificate paths/secrets. With nothing matched, it emits only its error
counter.

**How to check:**

1. Confirm the selectors against reality — check the exporter args (`--secrets-include-glob`,
   `--secrets-annotation-selector`, disk `--include-glob`) versus the actual secret annotations / file paths present.
2. Check the exporter logs for read/permission errors on the cert sources.

**How to fix:**

1. Correct the glob/annotation selector to match the real certificate secrets or file paths, and ensure the
   ServiceAccount/mounts grant read access to them.

**How to avoid this issue:**

Validate the selector against a known certificate at deploy time, and alert if `cert_exporter_error_total` rises while
expiry metrics are absent.

**Data to collect:**

* Exporter args and logs.
* The annotations/paths of the certificates it should match.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [Single metric available - cert_exporter_error_total (Issue #33)](https://github.com/joe-elliott/cert-exporter/issues/33)
* [joe-elliott/cert-exporter README](https://github.com/joe-elliott/cert-exporter)
<!-- markdownlint-enable line-length -->

## Pushgateway

### Metrics pushed by batch jobs vanish after Pushgateway restarts

**Symptoms:**

* After a Pushgateway pod restart, previously pushed metrics are gone.
* Dashboards for batch jobs show gaps aligned with the restart.

**Root cause:**

By default Pushgateway keeps metrics only in memory and does not persist across restarts. A pod restart (or crash) loses
everything unless `--persistence.file` is set and backed by a volume.

**How to check:**

1. Check whether persistence is enabled: inspect the Pushgateway args for `--persistence.file`
   (`kubectl get pod <pushgateway-pod> -n <namespace> -o yaml | grep persistence`).
2. Confirm whether the persistence path is on an `emptyDir` (lost on reschedule) rather than a PVC.

**How to fix:**

1. Set `spec.pushgateway.storage` on the `PlatformMonitoring` CR to a PVC spec. The operator then adds the
   `--persistence.file` and `--persistence.interval` flags and mounts the volume automatically — do not add the flags or
   the volume by hand.

**How to avoid this issue:**

Always set `spec.pushgateway.storage` so Pushgateway persists to durable storage; treat it as a cache that must survive
pod churn.

**Data to collect:**

* Pushgateway args and volume mounts.
* Restart timestamps versus the metric gap.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `controllers/pushgateway/manifest.go`, `api/v1/platformmonitoring_types.go`
* [Pushgateway - any way to save/restore metrics (Google Groups)](https://groups.google.com/g/prometheus-users/c/Ah4TsPpIT68)
* [prometheus/pushgateway README (--persistence.file)](https://github.com/prometheus/pushgateway/blob/master/README.md)
<!-- markdownlint-enable line-length -->

### Old batch-job metrics never disappear and alerts fire on stale values

**Symptoms:**

* A completed or one-off job's metrics remain visible indefinitely.
* Alerts fire on values that are hours or days old.

**Root cause:**

Pushgateway does not expire metrics — there is no TTL. Anything pushed stays "current" to Prometheus until it is
explicitly deleted via the DELETE API. Jobs that never clean up leave stale series behind.

**How to check:**

1. List what is held under each grouping key on the Pushgateway `/metrics` endpoint and read `push_time_seconds` to see
   how old each group is.

**How to fix:**

1. Have each job DELETE its group on completion or failure:
   `curl -X DELETE http://<pushgateway>:9091/metrics/job/<job>/<label>/<value>`.
2. For alerting, gate on `push_time_seconds` freshness so a stale push cannot keep firing an alert.

**How to avoid this issue:**

Push under stable grouping keys, delete on job exit, and use Pushgateway only for genuine batch or ephemeral jobs, not
for machine-level monitoring.

**Data to collect:**

* The Pushgateway `/metrics` group list with `push_time_seconds`.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [When metrics disappear on updates with Prometheus Pushgateway (Chris Siebenmann)](https://utcc.utoronto.ca/~cks/space/blog/sysadmin/PrometheusPushgatewayDropMetrics)
* [prometheus/pushgateway README (DELETE / no expiry)](https://github.com/prometheus/pushgateway/blob/master/README.md)
<!-- markdownlint-enable line-length -->

## Prometheus-adapter

### ServiceNotFound for the v1beta1.custom.metrics.k8s.io API service

**Symptoms:**

* Kubernetes API server cannot handle some requests; a namespace stays stuck after a delete command.
* Logs show:

  <!-- markdownlint-disable line-length -->
  ```text
  "Error while check hasAPI","error":"unable to retrieve the complete list of server APIs: custom.metrics.k8s.io/v1: the server is currently unable to handle the request"
  ```
  <!-- markdownlint-enable line-length -->

* The API service is not available:

  ```text
  v1beta1.custom.metrics.k8s.io   prometheus/prometheus-adapter   False (ServiceNotFound)   12s
  ```

**Root cause:**

`prometheus-adapter` is not available (the pod was removed or is crash-looping), but its `APIService` is still
registered. Because the adapter handles a cluster-wide aggregated API, its absence blocks any request routed to that API
— including deleting a namespace that carries the `kubernetes` finalizer. Bringing the adapter pod back is the whole of
the **prometheus-adapter pod is down or restarting** case; see also that case.

**How to check:**

1. Read the API service status; `False` in the `AVAILABLE` column confirms the adapter cannot serve requests:

   ```bash
   kubectl get apiservices
   ```

**How to fix:**

1. Preferred: restore the `prometheus-adapter` pod — run it again or give it more resources — so it serves the
   registered API.
2. **DANGEROUS — deleting the APIService unregisters the custom metrics API cluster-wide; HPAs relying on custom
   metrics stop working until the adapter re-registers it.** If the adapter is gone for good, remove the stale
   registration:

   ```bash
   kubectl delete apiservice v1beta1.custom.metrics.k8s.io
   ```

   This chart registers the APIService as `v1beta1.custom.metrics.k8s.io` (the custom-metrics API group has no `v1`), so
   match on that exact name.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`, `docs/user-guides/horizontal-autoscaling.md`
<!-- markdownlint-enable line-length -->

### FailedDiscoveryCheck for the v1beta1.custom.metrics.k8s.io API service

**Symptoms:**

* HPA by custom metrics does not work; HPA objects emit an event:

  ```text
  Warning  FailedGetPodsMetric  <invalid> (x63342 over 10d)  horizontal-pod-autoscaler
    unable to get metric tm_busyness: unable to fetch metrics from custom metrics API:
      the server could not find the metric <metric_name> for pods
  ```

* The API service shows `False (FailedDiscoveryCheck)`:

  ```text
  v1beta1.custom.metrics.k8s.io   monitoring/prometheus-adapter   False (FailedDiscoveryCheck)   85m
  ```

* Or the whole API is healthy but one custom metric never appears in it, while the adapter pod is healthy and other
  metrics work.

**Root cause:**

`prometheus-adapter` is available but discovered no rules, so it does not run the handler for the registered API. Either
the label selector for discovering `CustomScaleMetricRule` CRs is wrong, or there are no such CRs with a rule at all.
The adapter enables the handler only when it has at least one rule for `v1beta1.custom.metrics.k8s.io`. The same
selector-mismatch mechanism produces a narrower, single-rule case: when the API is otherwise healthy, one
`CustomScaleMetricRule` that does not match the operator's `customScaleMetricRulesSelector` is silently dropped while
the operator merges rules into the adapter configuration, so that one metric never reaches the custom metrics API even
though others work.

**How to check:**

1. Read the API service status for `FailedDiscoveryCheck`:

   ```bash
   kubectl get apiservices
   ```

2. Confirm at least one `CustomScaleMetricRule` exists and carries the labels the adapter's selector requires. For the
   single-missing-metric case, check the specific rule's labels against the operator's `customScaleMetricRulesSelector`
   and confirm the metric is absent from the API:

   ```bash
   kubectl -n <namespace> get customscalemetricrule <name> -o jsonpath='{.metadata.labels}{"\n"}'
   kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1 | tr ',' '\n' | grep <metric>
   ```

**How to fix:**

1. If the selector is too strict, either clear it or add the expected labels to the CR. Clear it with:

   ```yaml
   prometheusAdapter:
     customScaleMetricRulesSelector: []
   ```

   Or add the expected label to the CR, for example `app.kubernetes.io/component: monitoring`. The same fix resolves the
   single-missing-metric case: label the one dropped `CustomScaleMetricRule` so it satisfies the selector, or relax the
   selector, then wait for the next reconcile.
2. If there are no rules, create at least one `CustomScaleMetricRule` with at least one rule so the adapter enables the
   handler.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`, [controllers/customscalemetricrule_controller.go](https://github.com/Netcracker/qubership-prometheus-adapter-operator/blob/main/controllers/customscalemetricrule_controller.go)
<!-- markdownlint-enable line-length -->

### prometheus-adapter pod is down or restarting

**Symptoms:**

* The `prometheus-adapter` pod is down or restarts continuously.

**Root cause:**

There are three common causes: the adapter was OOMKilled for lack of resources; `Prometheus URL` points at the wrong
address; or the Prometheus it reads from is unavailable. While the adapter pod is down, its registered `APIService`
fails cluster-wide; see also **ServiceNotFound for the v1beta1.custom.metrics.k8s.io API service**.

**How to check:**

1. Read the pod's last state and events for `OOMKilled`:

   ```bash
   kubectl -n <namespace> describe pod <prometheus-adapter-pod>
   ```

2. Read the configured Prometheus URL and confirm the target Prometheus is reachable.

**How to fix:**

1. For OOM, raise CPU/memory limits in the deploy parameters or the `PrometheusAdapter` CR:

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

2. For a wrong link, set a full Prometheus URL with schema and port:

   ```yaml
   prometheusAdapter:
     prometheusUrl: http://prometheus.monitoring.svc:9090
   ```

3. If Prometheus itself is unavailable, fix Prometheus (see the Prometheus cases).

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### A second PrometheusAdapter custom resource is never reconciled

**Symptoms:**

* A newly created `PrometheusAdapter` CR is ignored; the operator log reads:

  ```text
  skip reconciliation: now reconcile <ns/name> (retry after 3m0s)
  ```

* Before the adapter is up, `CustomScaleMetricRule` reconciles are also deferred:

  ```text
  skip reconciliation: there is no reconciled prometheus-adapter (retry after 1m0s)
  ```

**Root cause:**

The operator manages a single `PrometheusAdapter` instance. Once one is active, any additional `PrometheusAdapter` CR is
permanently skipped. `CustomScaleMetricRule` reconciles are also deferred until a `PrometheusAdapter` is reconciled.

**How to check:**

1. Read the operator log for the `skip reconciliation` lines and note which instance it is reconciling.
2. List `PrometheusAdapter` CRs and confirm more than one exists:

   ```bash
   kubectl get -n <namespace> prometheusadapters.monitoring.netcracker.com
   ```

**How to fix:**

1. Keep exactly one `PrometheusAdapter` CR per operator; remove the extra CR and let the single instance reconcile.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-prometheus-adapter-operator](https://github.com/Netcracker/qubership-prometheus-adapter-operator) (not yet confirmed on a live install).
* Derived from source: [controllers/prometheusadapter_controller.go](https://github.com/Netcracker/qubership-prometheus-adapter-operator/blob/main/controllers/prometheusadapter_controller.go)
<!-- markdownlint-enable line-length -->

### Custom metrics occasionally disappear from the adapter

**Symptoms:**

* Custom metrics served by prometheus-adapter intermittently vanish and reappear.

**Root cause:**

`metricsRelistInterval` is smaller than the Prometheus scrape interval. The adapter only lists metrics that exist
between the current time and the last discovery query, so a relist interval shorter than the scrape interval causes
metrics to drop out between relists.

**How to check:**

1. Compare the configured `metricsRelistInterval` against the Prometheus scrape interval. A relist interval smaller than
   the scrape interval is the cause.

**How to fix:**

1. Set `metricsRelistInterval` equal to or larger than the Prometheus scrape interval, then redeploy.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-prometheus-adapter-operator](https://github.com/Netcracker/qubership-prometheus-adapter-operator) (not yet confirmed on a live install).
* Derived from source: [docs/install.md](https://github.com/Netcracker/qubership-prometheus-adapter-operator/blob/main/docs/install.md)
<!-- markdownlint-enable line-length -->

### Prometheus adapter replaces the in-built HPA resource metrics adapter

By default this deployment registers `prometheus-adapter` only for custom metrics (`v1beta1.custom.metrics.k8s.io`); it
does not take over the resource-metrics API, so the in-built provider (metrics-server) keeps serving CPU and memory. The
resource-metrics APIService (`v1beta1.metrics.k8s.io`) is registered only when `enableResourceMetrics` and
`APIService.resourceMetrics` are both set to `true` (and `global.privilegedRights` is on) — the subchart defaults both
to `false`. Enable those and `prometheus-adapter` replaces the in-built resource-metrics adapter, after which CPU- and
memory-based scaling flows through it. To keep that scaling working, a default configuration exposes container metrics
in two places — the aggregated ConfigMap for `prometheus-adapter` and a default `CustomScaleMetricRule` CR named
`kubelet-custom` — and the adapter cannot start with an empty configuration, so the default must be present in both.

Do not remove the default `kubelet-custom` CR or add another CR before it is in place: `prometheus-adapter` must never
run with an empty configuration. This is background for the FailedDiscoveryCheck case above (including its
single-dropped-rule scenario) rather than a failure in itself.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/user-guides/horizontal-autoscaling.md`
<!-- markdownlint-enable line-length -->

## Integrations

### graphite-remote-adapter regularly restarts with OOM

**Symptoms:**

* The graphite-remote-adapter pod restarts with an `OOMKilled` status.
* It can start after a restart, after being scaled up from 0, or move between environments and then begin
  OOM-restarting.

**Root cause:**

Prometheus RemoteWrite replays up to the last 2 hours of write-ahead-log data when the adapter comes back after being
down. That flood of points is buffered in the adapter's cache, which can grow past the memory limit and trigger the
OOMKiller.

**How to check:**

1. Confirm the pod's last state is `OOMKilled` and correlate the restart with the adapter having been down, scaled to 0,
   or newly started:

   ```bash
   kubectl -n <namespace> describe pod <graphite-remote-adapter-pod>
   ```

**How to fix:**

1. Increase resources (recommended): raise CPU to 1-2 vCPU and memory to 2-4 GB and redeploy:

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

2. Alternatively, decrease the cache size by tuning the Graphite write config (or the `graphite-config` ConfigMap), then
   restart the Graphite pod:

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

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### graphite-remote-adapter sends metrics with delay

**Symptoms:**

* graphite-remote-adapter sends metrics to Graphite / carbon-c-relay with a delay of 5 minutes to 1 hour.

**Root cause:**

Either the downstream (carbon-c-relay / Graphite) cannot process the full metric stream, or the adapter lacks resources
or is CPU-throttled and cannot convert and send fast enough.

**How to check:**

1. Confirm carbon-c-relay and Graphite can process the whole stream the adapter produces.
2. Enable and read the adapter's own metrics to see whether it is under-resourced or throttled:

   ```yaml
   graphite_remote_adapter:
     servicemonitor:
       install: true
     grafanaDashboard: true
   ```

**How to fix:**

1. If the adapter is under-resourced, scale it out. Update the replica count and redeploy, or scale in place:

   ```bash
   kubectl --namespace=monitoring scale deployment --replicas=2 \
     --selector="app.kubernetes.io/component=graphite-remote-adapter"
   ```

2. If the adapter is heavily throttled, set `GOMAXPROCS` to match the CPU limit (for a 2-core limit, `GOMAXPROCS: "2"`).
   Since release `0.70.0` the adapter sets `GOMAXPROCS` automatically from the pod limits, so skip this on that version
   or later.

**Sources:**

<!-- markdownlint-disable line-length -->
* Derived from source: `docs/troubleshooting.superseded.md`
<!-- markdownlint-enable line-length -->

### graphite-remote-adapter silently drops metrics with no carbon address set

**Symptoms:**

* Metrics written through graphite-remote-adapter never arrive at Graphite, yet RemoteWrite reports success (HTTP 200).
* The adapter response body reads:

  ```text
  Skipped: Not set carbon address.
  ```

**Root cause:**

When the carbon address is empty, the write path returns HTTP 200 with the `Skipped: Not set carbon address.` body and a
nil error, so nothing is written and Prometheus sees a successful remote write.

**How to check:**

1. Confirm the configured carbon address is empty in the adapter configuration.
2. Capture the adapter response body for a write and look for `Skipped: Not set carbon address.`.

**How to fix:**

1. Set a valid carbon address in the adapter configuration and redeploy so writes are forwarded.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-graphite-remote-adapter](https://github.com/Netcracker/qubership-graphite-remote-adapter) (not yet confirmed on a live install).
* Derived from source: [client/graphite/write.go](https://github.com/Netcracker/qubership-graphite-remote-adapter/blob/main/client/graphite/write.go)
<!-- markdownlint-enable line-length -->

### graphite-remote-adapter rejects the whole config on an unknown field

**Symptoms:**

* graphite-remote-adapter fails to load its configuration after an edit to `additionalGraphiteConfig`, logging:

  ```text
  unknown fields in <context>: <field-list>
  ```

**Root cause:**

Configuration parsing is strict: any unknown field — typically a typo — aborts the entire configuration load, not
just the offending section.

**How to check:**

1. Read the adapter log for the `unknown fields in` error and note the context (graphite config, write config, read
   config, or rule) and the field list.

**How to fix:**

1. Correct or remove the unknown field named in the error so the configuration parses, then reload or redeploy.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from sibling repository [qubership-graphite-remote-adapter](https://github.com/Netcracker/qubership-graphite-remote-adapter) (not yet confirmed on a live install).
* Derived from source: [utils/config.go](https://github.com/Netcracker/qubership-graphite-remote-adapter/blob/main/utils/config.go)
<!-- markdownlint-enable line-length -->

## Promxy

### sum()/count() through Promxy return doubled or un-aggregated values across backends

**Symptoms:**

* Aggregations like `sum(up)` or `count(...)` return larger-than-expected results through Promxy.
* The same query against a single backend returns the correct value.

**Root cause:**

Promxy de-duplicates within a server group by labelset. When the same series exists across multiple server groups (or
server groups carry distinguishing labels), Promxy treats them as distinct and does not aggregate across them, so
HA/sharded backends produce inflated aggregation results.

**How to check:**

1. Compare the aggregation result from Promxy against the result from one backend directly to confirm the doubling.
2. Inspect the Promxy `server_groups` config — check whether the same targets appear in multiple groups or whether
   groups add labels (`labels:`) that split otherwise-identical series.

**How to fix:**

1. Place HA replicas of the same data in one server group (with `anti_affinity` set to your scrape interval) so Promxy
   de-duplicates them, and reserve separate server groups for genuinely distinct data.

**How to avoid this issue:**

Model server groups so that replicas of identical data share a group and distinct shards get distinct groups; avoid
adding group labels that fragment identical series.

**Data to collect:**

* The Promxy result vs. single-backend result for the same query.
* The Promxy `server_groups` configuration.

**Sources:**

<!-- markdownlint-disable line-length -->
* Compiled from external research (upstream issues and vendor documentation); not confirmed on a live installation.
* [promxy doesn't correctly aggregate across multiple Promethei (Issue #260)](https://github.com/jacksontj/promxy/issues/260)
* [Promxy do not aggregate sum and count functions (Issue #368)](https://github.com/jacksontj/promxy/issues/368)
<!-- markdownlint-enable line-length -->
