# Monitoring Pack Quickstart

This quickstart covers two deliverables:

1. **Operators-only baseline** – install Monitoring Operator, VictoriaMetrics Operator, and Grafana Operator with
nothing except the mandatory exporters.
2. **Pack-one** – extend that baseline with the minimal OPS bundle (vmagent, vmalert, ServiceMonitors, remote write).

Follow the sections below in order. The steps assume you already cloned this repository and have `kubectl` and `helm`
configured against your target cluster.

> **Namespaces used**
>
> - `monitoring` – all operators (Monitoring Operator, VictoriaMetrics Operator, Grafana Operator) and pack-one components.
> - `vmsingle-standalone` – a standalone VictoriaMetrics `vmsingle` instance for functional checks.

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

This step covers the pack-one add-on. The chart delivers:

- `vmagent` configured with Remote Write to the test `vmsingle`
- `vmalert` plus default ServiceMonitors for the mandatory exporters
- `vmauth` for authentication and routing (optional)
- `vmuser` CRD for user credentials and routing configuration (optional)
- Ingress resources for external access (optional)
- RBAC required by VictoriaMetrics Operator

Install it into the same `monitoring` namespace:

```bash
helm upgrade --install pack-one docs/examples/deploy-parameters/monitoring-pack/pack-one \
  --namespace monitoring \
  -f docs/examples/deploy-parameters/monitoring-pack/pack-one/values.yaml
```

### 4.1. Configure VMAuth and VMUser (optional)

By default, VMAuth and VMUser are installed with basic authentication (`admin/admin`). VMAuth automatically discovers VMUser resources using the label selector `app.kubernetes.io/name: vmuser`. Routing is automatically configured to route requests to installed components (vmagent, vmsingle) based on path patterns.

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

2. **Customize routing** by specifying `targetRefs` in `vmUser.spec.targetRefs`. If not specified, targetRefs are automatically generated for installed components:
   - **vmagent**: routes `/config.*`, `/target.*`, `/service-discovery.*`, `/static.*`, `/api/v1/write`, `/api/v1/import.*`, `/api/v1/target.*`
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

2. **Update DNS** to point your hosts to the Ingress controller's IP address.

3. **Access services**:
   - `vmauth`: `http://vmauth.example.com` (or `https://` if TLS is configured)
   - `vmagent`: `http://vmagent.example.com` (or `https://` if TLS is configured)
   - `vmsingle`: `http://vmsingle.example.com` (or `https://` if TLS is configured)

   > **Note**: If VMAuth is enabled, access vmagent and vmsingle through VMAuth using the configured credentials.

## 5. Access services

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
kubectl port-forward -n monitoring svc/vmauth-vmauth 18427:8427 &

# VMAgent
kubectl port-forward -n monitoring svc/vmagent-vmagent 18428:8429 &

# VMSingle
kubectl port-forward -n vmsingle-standalone svc/vmsingle-k8s 18429:8429 &
```

Access the UIs/APIs through a browser or API client:

- `vmauth`: [http://127.0.0.1:18427](http://127.0.0.1:18427) (if enabled, use `admin/admin` for authentication)
  - VMUI through VMAuth: [http://127.0.0.1:18427/vmui/](http://127.0.0.1:18427/vmui/)
  - vmagent targets through VMAuth: [http://127.0.0.1:18427/targets](http://127.0.0.1:18427/targets)
- `vmagent`: [http://127.0.0.1:18428](http://127.0.0.1:18428)
- `vmsingle`: [http://127.0.0.1:18429](http://127.0.0.1:18429)

## 6. Validate the data flow

1. **vmagent targets** – open the vmagent UI (via Ingress or port-forward) → `/targets` and check that all
mandatory exporter targets are in the `UP` state.

   ![vmagent targets](./images/vmagent-example.png)

2. **vmsingle query** – open the vmsingle UI (via Ingress or port-forward) → `/vmui/` → `Query`, type `up`, run the
query, and confirm that vmagent's scraped metrics arrive in vmsingle (results should include the exporters you enabled).

   ![vmsingle query](./images/vmsingle-example.png)

3. **VMAuth routing** (if enabled) – access vmagent or vmsingle through VMAuth using the configured credentials (`admin/admin` by default).
   - Access VMUI: `http://vmauth.example.com/vmui/` (or via port-forward: `http://127.0.0.1:18427/vmui/`)
   - Access vmagent targets: `http://vmauth.example.com/targets`
   - Access vmsingle query: `http://vmauth.example.com/vmui/query`
   
   The VMAuth service automatically routes requests to the appropriate backend (vmagent or vmsingle) based on the path patterns defined in VMUser's targetRefs. The root path `/` is included to support VMUI access.

---

Outcome:

1. Baseline operators installed with only mandatory exporters (task “operators-only baseline”).
2. Pack-one deployed on top of that baseline, sending metrics to the test `vmsingle` (task “pack-one OPS bundle”).
