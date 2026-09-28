# Infrastructure component troubleshooting


## Shared discovery

### Monitor excluded

**Symptoms:**

- A component `ServiceMonitor` or `PodMonitor` exists in the component namespace.
- The namespace label does not match the Prometheus or VMAgent `serviceMonitorNamespaceSelector`.
- The selector allowlist includes namespaces such as `monitoring` and `cqp` and omits the component namespace.

**Root cause:**

Monitoring Operator installation values restrict which namespaces VMAgent or Prometheus may discover. A valid
component monitor outside that scope is never scraped.

**Applies to:** any supported component after its scrape resource exists.

**Failed layer:** `monitor excluded`

**How to check:**

1. Confirm from the ticket that the component scrape resource exists. If the operator still needs to collect it,
   they can run:

   ```bash
   kubectl --context <context> -n <component-namespace> get servicemonitors,podmonitors
   kubectl --context <context> -n <component-namespace> get servicemonitor <name> -o yaml
   ```

2. Locate the `PlatformMonitoring` resource and read its VMAgent or Prometheus namespace and label selectors from
   pasted YAML, or ask the operator to collect:

   ```bash
   kubectl --context <context> get platformmonitoring -A
   kubectl --context <context> -n <monitoring-namespace> get platformmonitoring <name> -o yaml
   ```

3. Compare those selectors with the monitor namespace labels. Stop here when they do not match.

If the ticket already quotes the selector and the monitor namespace, that is enough. Do not enable a component
exporter flag to fix this layer. Do not run `kubectl` from this skill.

**How to fix:**

1. Read `victoriametrics.vmAgent.serviceMonitorNamespaceSelector` (or the Prometheus equivalent) in the
   Monitoring Operator installation values.
2. Ask the Monitoring Operator owners to extend that selector, or to apply a namespace label the selector
   already requires.
3. After they apply the values, re-read the PlatformMonitoring spec and confirm it selects the component
   namespace.

**Owner:** Monitoring Operator installation values

**Risk:** Expanding the namespace selector causes VMAgent to discover ServiceMonitors in every extra matching
namespace. Confirm that discovery scope before applying the values change.

**Provenance:** Monitoring Operator docs `docs/monitoring-configuration/limits-metric-collection.md` and
`docs/installation/components/victoriametrics-stack/vmagent.md`.


### Target unavailable

**Symptoms:**

- The expected `ServiceMonitor` or `PodMonitor` exists and is in discovery scope.
- The selected Service has no ready endpoints, or the endpoint port does not match the monitor.

**Root cause:**

The monitor does not select a ready Service port. The scrape target never becomes healthy.

**Applies to:** any supported component after the monitor is discovered.

**Failed layer:** `target unavailable`

**How to check:**

Read the pasted Service and EndpointSlice YAML. If the operator still needs to collect it, they can run:

```bash
kubectl --context <context> -n <component-namespace> get service <name> -o yaml
kubectl --context <context> -n <component-namespace> get endpointslice -l kubernetes.io/service-name=<name> -o yaml
```

Stop when the Service selector, endpoint port, or ready endpoints do not match the monitor. A scrape timeout on a
ready target is the component collection-timeout case (pgskipper `Collection scrape timeout`, MongoDB `MongoDB
collection scrape timeout`, or ClickHouse `ClickHouse collection scrape timeout`), not this one. Do not run `kubectl`
from this skill.

**How to fix:**

1. Ask the component team to align the Service selector and the named metrics port with the monitor.
2. Do not change Monitoring Operator selectors for a target that is not ready.

**Owner:** component chart or workload

**Risk:** None — every recommended step is read-only

**Provenance:** none beyond the live Service and EndpointSlice YAML.


### Access failure

**Symptoms:**

- The pasted output shows that the named Kubernetes context, namespace, resource type, or permission is unavailable.
- Example: `error: context "production-eu" does not exist`.

**Root cause:**

The diagnosis cannot inspect the named context. Continue from ticket attachments or pasted non-secret evidence.

**Applies to:** after the ticket quotes a named `--context` and namespace, before cluster resources can be read.

**Failed layer:** `access failure`

**How to check:**

Quote the command and the error verbatim from the ticket. Do not run `kubectl`. Do not run
`kubectl config use-context`. Do not retry against another context.

**How to fix:**

Request one of:

- The ticket description and non-secret attachments.
- The complete image reference from the component workload.
- The non-secret workload YAML showing `.spec.template.spec.containers[*].image`.
- The same read-only `kubectl --context <name>` output from a kubeconfig that already contains that context.

Do not ask for Secret data, credentials, or kubeconfig content.

**Owner:** none. Diagnosis stopped before a metric-path layer.

**Risk:** None — every recommended step is read-only

**Provenance:** none. Load no component case until evidence from the ticket is available.
