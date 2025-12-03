# Monitoring Pack Quickstart

This quickstart covers three deliverables:

1. **Operators-only baseline** – install Monitoring Operator, VictoriaMetrics Operator, and Grafana Operator with
nothing except the mandatory exporters.
2. **Pack-one** – extend that baseline with the minimal OPS bundle (vmagent, vmalert, ServiceMonitors, remote write)
   plus shared exporter configuration.
3. **Pack-two** – provision a DEV-focused monitoring bundle (dedicated VMSingle/VMAgent/VMAlert, Grafana objects, and
   ingress) that reuses the exporters from pack-one via label selectors.

Follow the sections below in order. The steps assume you already cloned this repository and have `kubectl` and `helm`
configured against your target cluster.

> **Namespaces used**
>
> - `monitoring` – all operators (Monitoring Operator, VictoriaMetrics Operator, Grafana Operator);
> - `pack-one` – pack-one components (VMAgent, VMAlert, VMAuth, VMUser);
> - `pack-two` – pack-two components (DEV VMSingle/VMAgent/VMAlert/VMAuth, Grafana resources);
> - `vmsingle-standalone` – a standalone VictoriaMetrics `vmsingle` instance for functional checks;

## 1. Install Grafana Operator CRDs

Grafana Operator CRDs must exist before Helm installs any resources that refer to them:

```bash
kubectl apply -f https://raw.githubusercontent.com/grafana/grafana-operator/v4.10.0/deploy/manifests/v4.10.0/crds.yaml
```

## 2. Deploy the baseline Monitoring Operator

Install the Monitoring Operator with only the mandatory exporters enabled (baseline scenario):

```bash
helm upgrade --install qubership-monitoring-operator \
  charts/qubership-monitoring-operator \
  --namespace monitoring \
  --create-namespace \
  -f docs/examples/deploy-parameters/monitoring-pack/operators-only-baseline-values.yaml
```

This chart deploys Monitoring Operator, VictoriaMetrics Operator, and Grafana Operator without any dashboards,
datasources, or VictoriaMetrics clusters.

## 3. Deploy a standalone vmsingle (optional test target)

Create an isolated namespace and apply the provided manifest to run a minimal VictoriaMetrics `vmsingle`. This instance
is useful for validating remote write from `vmagent`:

```bash
kubectl create namespace vmsingle-standalone
kubectl apply -f docs/examples/deploy-parameters/monitoring-pack/tests/vmsingle-direct.yaml
```

## 4. Install pack-one (OPS monitoring bundle)

Pack-one acts as the operational baseline: it wires the mandatory exporters to a managed VMAgent/VMSingle pair, ships a
starter alerting pipeline, and exposes the endpoints through VMAuth/Ingress so OPS teams can observe the platform even
before DEV-focused features (like pack-two dashboards) are installed. The chart delivers:

- `vmagent` configured with Remote Write to the test `vmsingle`
- `vmalert` plus default ServiceMonitors for the mandatory exporters
- `vmauth` for authentication and routing (optional)
- `vmuser` CRD for user credentials and routing configuration (optional)
- Ingress resources for external access (optional)
- RBAC required by VictoriaMetrics Operator

Install it in `pack-one` namespace:

```bash
helm upgrade --install pack-one docs/examples/deploy-parameters/monitoring-pack/pack-one \
  --namespace pack-one \
  --create-namespace \
  -f docs/examples/deploy-parameters/monitoring-pack/pack-one/values.yaml
```

> **Shared exporters label**: every ServiceMonitor rendered by pack-one automatically receives the
> `monitoring-pack: "one"` label. Override `monitoringPackLabel` in `values.yaml` if you need a different value.
> Downstream bundles (for example, pack-two) rely on this label to reuse exporters.

### 4.1. Configure VMAuth and VMUser (optional)

By default, VMAuth and VMUser are installed with basic authentication (`admin/admin`). VMAuth automatically discovers
VMUser resources using the label selector `app.kubernetes.io/name: vmuser`. Routing is automatically configured to route
requests to installed components (vmagent, vmsingle) based on path patterns.

To customize:

1. **Update credentials** in `values.yaml`:
   ```yaml
   vmUser:
     install: true
     spec:
       username: your-username
       password: your-password
       # Or use passwordRef for secret-based authentication
       # passwordRef:
       #   name: vmuser-secret
       #   key: password
   ```

2. **Customize routing** by specifying `targetRefs` in `vmUser.spec.targetRefs`. If not specified, targetRefs are
automatically generated for installed components:
   - **vmagent**: routes `/config.*`, `/target.*`, `/service-discovery.*`, `/static.*`, `/api/v1/write`,
   `/api/v1/import.*`, `/api/v1/target.*`
   - **vmsingle**: routes `/` (root path for VMUI), `/vmui.*`, `/graph.*`, `/api/v1/*`, `/prometheus/*`, etc.

3. **VMAuth configuration** (if needed):
   ```yaml
   vmAuth:
     install: true
     spec:
       selectAllByDefault: false  # Default: false - uses userSelector to find VMUser
       # userSelector: {}  # Default: matches app.kubernetes.io/name: vmuser
       # userNamespaceSelector: {}  # Default: empty (searches in all namespaces)
   ```

### 4.2. Configure Ingress (optional)

To enable Ingress resources instead of port-forwarding:

1. **Control Ingress installation** via `ingress.install` (default `true`). Sample configuration:
   ```yaml
   ingress:
     install: true
     vmAuth:
       host: vmauth.example.com
       ingressClassName: traefik # default controller in the quickstart cluster
       servicePort: 8427
       # Optional: TLS configuration
       # tlsSecretName: vmauth-tls
     vmAgent:
       host: vmagent.example.com
       ingressClassName: traefik
       servicePort: 8429
     vmSingle:
       host: vmsingle.example.com
       ingressClassName: traefik
       servicePort: 8429
   ```

3. **Access services**:
   - `vmauth`: `http://vmauth.example.com` (or `https://` if TLS is configured)
   - `vmagent`: `http://vmagent.example.com` (or `https://` if TLS is configured)
   - `vmsingle`: `http://vmsingle.example.com` (or `https://` if TLS is configured)

   > **Note**: If VMAuth is enabled, access vmagent and vmsingle through VMAuth using the configured credentials.

## 5. Install pack-two (DEV monitoring bundle)

Pack-two targets development teams that need their own VictoriaMetrics + Grafana stack while still scraping the shared
exporters delivered by pack-one. The chart installs:

1. **VMSingle** – a DEV-only VictoriaMetrics cluster (`vmSingle`)
2. **VMAgent** – collects metrics from both `monitoring-pack=one` and `monitoring-pack=two` ServiceMonitors
3. **VMAlert** – rules and alert delivery scoped to DEV workloads
4. **VMAuth / VMUser** (optional) – authentication and routing facade for DEV access
5. **ServiceMonitors** – pack-two specific exporters (each labeled `monitoring-pack: "two"`)
6. **Grafana objects** – datasource, example dashboards, and alerting rules bound to the DEV VMSingle

### 5.1 Reusing pack-one exporters

1. Ensure pack-one ServiceMonitors carry `monitoring-pack: "one"` (this is automatic unless you override
   `monitoringPackLabel`).
2. Pack-two’s VMAgent is rendered with `serviceScrapeSelector.matchExpressions` that includes both `"one"` and `"two"`,
   so it automatically discovers the shared exporters alongside its own.
3. Each pack still keeps an isolated storage/alerting stack—only the scrape configuration is shared.

### 5.2 Install pack-two

> **Note**: Pack-two assumes VictoriaMetrics Operator is already running (from baseline/pack-one) in the `monitoring`
> namespace. If your operator lives elsewhere, update `vmOperatorNamespace` in
> `docs/examples/deploy-parameters/monitoring-pack/pack-two/values.yaml` before installing.

```bash
helm upgrade --install pack-two docs/examples/deploy-parameters/monitoring-pack/pack-two \
  --namespace pack-two \
  --create-namespace \
  -f docs/examples/deploy-parameters/monitoring-pack/pack-two/values.yaml
```

### 5.3 Configure VMAuth and VMUser (optional)

By default pack-two deploys `vmuser-dev` with `admin/admin` credentials and labels that limit VMAuth discovery to pack-two
users only. Override these credentials (or switch to secrets/tokens) before going to production:

```yaml
vmUser:
  install: true
  spec:
    username: your-username
    passwordRef:
      name: vmuser-secret
      key: password
    # or use bearerToken / tokenRef for API access tokens
vmAuth:
  spec:
    # Optional: broaden discovery to multiple monitoring packs
    userSelector:
      matchLabels:
        app.kubernetes.io/name: vmuser
```

### 5.4 Configure Ingress for pack-two (optional)

```yaml
ingress:
  install: true
  vmAuth:
    host: vmauth-dev.example.com
    ingressClassName: traefik
    servicePort: 8427
  vmAgent:
    host: vmagent-dev.example.com
    ingressClassName: traefik
    servicePort: 8429
  vmSingle:
    host: vmsingle-dev.example.com
    ingressClassName: traefik
    servicePort: 8429
```

### 5.5 Validate pack-two

1. `kubectl get vmsingle,vmagent,vmalert -n pack-two`
2. Open the pack-two VMAgent UI (`http://vmagent-dev.example.com/targets`) and ensure exporter targets from both packs
   are in the `UP` state.
3. Open the pack-two VMSingle UI (`http://vmsingle-dev.example.com/vmui/`) and run `up{monitoring-pack="two"}`.

### 5.6 Grafana Architecture and Dashboards

#### How Grafana Works

The Grafana deployment follows a specific architecture pattern:

1. **Baseline deploys Grafana CR** – The Grafana CR is created in the `monitoring` namespace through PlatformMonitoring CR. This is required because Grafana Operator watches for Grafana CRs only in the namespace where it runs (`monitoring`).

2. **Grafana Operator scans multiple namespaces** – The operator is configured with `--namespaces="monitoring,pack-two"` to scan for dashboards and datasources in both namespaces.

3. **Pack-two creates dashboards** – Dashboards are created in the `pack-two` namespace with the label `monitoring-pack: "two"`.

4. **Grafana finds dashboards via selectors** – The Grafana CR from baseline uses `dashboardLabelSelector` and `dashboardNamespaceSelector` to discover dashboards in `pack-two` namespace.

5. **Dashboards connect to Grafana** – Each dashboard uses `instanceSelector` with `app.kubernetes.io/name: grafana` to connect to the Grafana instance from baseline.

6. **Datasource in pack-two** – The GrafanaDataSource is created in `pack-two` namespace because it connects to the VMSingle instance that exists in `pack-two`. Baseline doesn't have VMSingle, so the datasource must be where the metrics storage is located.

   > **⚠️ Warning: Datasource separation limitation**
   >
   > Currently, separating datasources between baseline and pack-two is not fully supported due to limitations in the Grafana Operator version. The Grafana Operator may not discover GrafanaDataSource resources in namespaces other than where it's deployed (`monitoring`), even with `instanceSelector` configured. This is a known limitation documented in [grafana-operator issue #304](https://github.com/grafana/grafana-operator/issues/304).
   >
   > For more details, see [GrafanaDataSource documentation](https://github.com/Netcracker/qubership-monitoring-operator/blob/main/docs/configuration.md?plain=1#L1385).
   >
   > **TBD:** Upgrade Grafana Operator to a newer version that fully supports cross-namespace datasource discovery with `instanceSelector`.

#### Mandatory Exporter Dashboards

Pack-two includes dashboards for all mandatory exporters:

- **nodeExporter** → `kubernetes-nodes-resources` dashboard (uses `node_*` metrics)
- **kubeStateMetrics** → `kubernetes-cluster-overview` dashboard (uses `kube_*` metrics)
- **cloudEventsExporter** → `cloud-events-exporter` dashboard
- **versionExporter** → `version-exporter` dashboard

These dashboards are located in `pack-two/templates/dashboards/` and can be enabled/disabled via `values.yaml`:

```yaml
grafana:
  dashboards:
    install: true
    mandatoryExporters:
      nodeExporter:
        enabled: true
      kubeStateMetrics:
        enabled: true
      cloudEventsExporter:
        enabled: true
      versionExporter:
        enabled: true
```

#### Adding Custom Dashboards

To add custom dashboards to pack-two:

1. **Create a dashboard template file** in `pack-two/templates/dashboards/`:

```yaml
{{- if and .Values.grafana.dashboards.install .Values.grafana.dashboards.custom.myDashboard.enabled }}
apiVersion: integreatly.org/v1alpha1
kind: GrafanaDashboard
metadata:
  name: {{ include "monitoring-pack-two.fullname" . }}-my-dashboard
  namespace: {{ include "monitoring-pack-two.namespace" . }}
  labels:
    app.kubernetes.io/component: monitoring
    monitoring-pack: "two"
    {{- $labelCtx := dict "commonLabels" .Values.commonLabels "componentLabels" dict }}
    {{- include "monitoring-pack-two.metadata.labels" $labelCtx | nindent 4 }}
spec:
  instanceSelector:
    matchLabels:
      app.kubernetes.io/name: grafana
  json: >
    {
      "title": "My Custom Dashboard",
      "panels": [...],
      ...
    }
{{- end }}
```

2. **Add configuration to values.yaml**:

```yaml
grafana:
  dashboards:
    install: true
    custom:
      - name: my-dashboard
        enabled: true
```

   Or if you prefer object-based structure, you can modify `values.yaml` to use an object instead of an array:

```yaml
grafana:
  dashboards:
    install: true
    custom:
      myDashboard:
        enabled: true
```

   Then update your template condition to match: `{{- if .Values.grafana.dashboards.custom.myDashboard.enabled }}`

3. **Important points**:
   - Dashboard must be in `pack-two` namespace
   - Must have label `monitoring-pack: "two"`
   - Must use `instanceSelector` with `app.kubernetes.io/name: grafana`
   - Use `${datasource}` variable in JSON for datasource references (will be resolved to pack-two datasource)

## 6. Access services

### Option A: Using Ingress (recommended)

If Ingress is enabled, access services through the configured hosts:

- **VMAuth**: `http://vmauth.example.com` (authenticate with credentials from `vmUser.spec`, default: `admin/admin`)
  - Access VMUI through VMAuth: `http://vmauth.example.com/vmui/`
  - Access vmagent through VMAuth: `http://vmauth.example.com/targets`
- **VMAgent**: `http://vmagent.example.com` (or through VMAuth if enabled)
- **VMSingle**: `http://vmsingle.example.com` (or through VMAuth if enabled)

### Option B: Using port-forward (for local troubleshooting)

If Ingress is not configured, expose services locally:

```bash
# VMAuth (if enabled)
kubectl port-forward -n pack-one svc/vmauth-vmauth 18427:8427 &

# VMAgent
kubectl port-forward -n pack-one svc/vmagent-vmagent 18428:8429 &

# VMSingle
kubectl port-forward -n vmsingle-standalone svc/vmsingle-k8s 18429:8429 &
```

Access the UIs/APIs through a browser or API client:

- `vmauth`: [http://127.0.0.1:18427](http://127.0.0.1:18427) (if enabled, use `admin/admin` for authentication)
  - VMUI through VMAuth: [http://127.0.0.1:18427/vmui/](http://127.0.0.1:18427/vmui/)
  - vmagent targets through VMAuth: [http://127.0.0.1:18427/targets](http://127.0.0.1:18427/targets)
- `vmagent`: [http://127.0.0.1:18428](http://127.0.0.1:18428)
- `vmsingle`: [http://127.0.0.1:18429](http://127.0.0.1:18429)

## 7. Validate the data flow

1. **vmagent targets** – open the vmagent UI (via Ingress or port-forward) → `/targets` and check that all
mandatory exporter targets are in the `UP` state.

   ![vmagent targets](./images/vmagent-example.png)

2. **vmsingle query** – open the vmsingle UI (via Ingress or port-forward) → `/vmui/` → `Query`, type `up`, run the
query, and confirm that vmagent's scraped metrics arrive in vmsingle (results should include the exporters you enabled).

   ![vmsingle query](./images/vmsingle-example.png)

3. **VMAuth routing** (if enabled) – access vmagent or vmsingle through VMAuth using the configured credentials
(`admin/admin` by default).
   - Access VMUI: `http://vmauth.example.com/vmui/` (or via port-forward: `http://127.0.0.1:18427/vmui/`)
   - Access vmagent targets: `http://vmauth.example.com/targets`
   - Access vmsingle query: `http://vmauth.example.com/vmui/query`
   
The VMAuth service automatically routes requests to the appropriate backend (vmagent or vmsingle) based on the path
patterns defined in VMUser's targetRefs. The root path `/` is included to support VMUI access.

## TODO

> **Note:** Current development is paused pending Grafana Operator upgrade. The following items need to be addressed:

- [ ] Check for redundant RBAC resources in pack-one/pack-two (pack-two RBAC was auto-generated)
- [ ] Reconfigure Ingress so that VMAgent and VMSingle route through VMAuth
- [ ] ServiceMonitors are deployed in the wrong namespace (`monitoring` instead of `pack-one`)
- [ ] Add dashboard examples and ability to override them
- [ ] Add recording rules examples and ability to override them
- [ ] Add alert examples and ability to override them
- [ ] **Q:** When creating pack-two, it also monitors pack-one. Is this correct? Why does this happen? Similar remote-write? (Need to remove selector "one" or "two" in agent)
