# Infrastructure component troubleshooting


## mongodb-operator

### Prometheus exporter is not installed

**Symptoms:**

- MongoDB 7 Prometheus or MongoDB 8 Prometheus dashboard panels have no data.
- Deployment and Service `mongodb-prometheus-exporter` are absent.

**Root cause:**

The mongodb-services chart does not install the exporter workload, so no scrape target exists. The documented
enablement is `prometheusExporter.install`, with helper fallback `MONITORING_ENABLED` (default true when both are
unset). This is a missing workload, not a scrape failure.

**Applies to:** mongodb-operator when MongoDB exporter metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `mongodb-prometheus-exporter` is absent.
2. Name `prometheusExporter.install` or `MONITORING_ENABLED` only when effective or rendered values show monitoring is
   off.
3. If the expected dashboard or series is unknown, request that title or series name, then still stop while the
   workload is missing.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `prometheusExporter.install: true` in the mongodb-services deployment values when that key is present in the
   evidence.
2. Re-render or upgrade the mongodb-services release using the component team's normal deployment process.
3. Confirm that Deployment and Service `mongodb-prometheus-exporter` exist, with labels
   `microservice: mongodb-prometheus-exporter`.

Do not change Monitoring Operator discovery selectors until the exporter workload exists.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Enabling the exporter deploys `mongodb-prometheus-exporter` and starts metric queries against MongoDB.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/cr_and_operator.yaml`
- `services/service/charts/helm/mongodb-services/templates/_helper.tpl`
- `services/service/pkg/prometheus_exporter/templates.go`
- `docs/public/installation_guide.md`


### Backup exporter monitor is absent

**Symptoms:**

- MongoDB backup dashboard panels such as Last Backup Status or Last Backup Size have no data.
- ServiceMonitor `mongo-backup-exporter-service-monitor` is absent.

**Root cause:**

The backup scrape resource is rendered only when backup is installed and monitoring is on. The documented gate is
`backup.install`. This is a missing backup exporter scrape, not the main MongoDB exporter.

**Applies to:** mongodb-operator when backup exporter metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `mongo-backup-exporter-service-monitor` is absent.
2. Name `backup.install` only when effective or rendered values show it is false, or the rendered templates omit that
   ServiceMonitor.

If those artifacts are missing, request them as Data required. Do not invent the key from a version label.

**How to fix:**

1. Set `backup.install: true` in the mongodb-services deployment values when that key is present in the evidence.
2. Re-render or upgrade the mongodb-services release using the component team's normal deployment process.
3. Confirm that `mongo-backup-exporter-service-monitor` selects `microservice: mongodb-backup-daemon` and path
   `/health/prometheus`.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Enabling backup creates `mongodb-backup-daemon` and an additional scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/values.yaml`
- `docs/public/grafana_dashboard.md`


### DBaaS adapter monitor is absent

**Symptoms:**

- MongoDB DBaaS Adapter dashboard row has no data.
- ServiceMonitor `mongo-dbaas-exporter-service-monitor` is absent.

**Root cause:**

The DBaaS adapter scrape resource is rendered only when the adapter is installed and monitoring is on. The documented
gates are `dbaas.install` and `DBAAS_ENABLED`.

**Applies to:** mongodb-operator when DBaaS adapter metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `mongo-dbaas-exporter-service-monitor` is absent.
2. Name `dbaas.install` or `DBAAS_ENABLED` only when effective or rendered values show the adapter is off.

If those artifacts are missing, request them as Data required.

**How to fix:**

1. Set `dbaas.install: true` in the mongodb-services deployment values when that key is present in the evidence.
2. Re-render or upgrade the mongodb-services release using the component team's normal deployment process.
3. Confirm that `mongo-dbaas-exporter-service-monitor` selects the DBaaS adapter and path `/metrics`.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Enabling the DBaaS adapter adds a scrape target for `dbaas-mongo-adapter`.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/templates/_helper.tpl`
- `docs/public/grafana_dashboard.md`


### MongoDB collection scrape timeout

**Symptoms:**

- ServiceMonitor `mongodb-prometheus-exporter-service-monitor` exists and is discovered, but the scrape ends before the
  exporter finishes.
- The target is down or partial.

**Root cause:**

This case is the chart `prometheusExporter.collectionScrapeTimeout` on the MongoDB exporter ServiceMonitor endpoint. It
is not the shared-discovery target-unavailable diagnosis (Service selector, endpoint port, or ready endpoints).

**Applies to:** mongodb-operator when a MongoDB exporter ServiceMonitor is discovered. Image and chart versions are
supporting facts.

**Failed layer:** `target unavailable`

**How to check:**

1. Confirm the expected ServiceMonitor exists and is in discovery scope.
2. Confirm the scrape timeout on that monitor from the live ServiceMonitor YAML or effective values
   (`prometheusExporter.collectionScrapeTimeout`).
3. Confirm the scrape ends before gathering finishes (target down or partial), not that endpoints are missing.

**How to fix:**

1. Raise the ServiceMonitor `scrapeTimeout` that the evidence named, then upgrade.
2. Do not use `prometheusExporter.collectionInterval` as the fix.
3. Do not change Monitoring Operator selectors for a scrape that already discovers the monitor.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Raising scrape timeout keeps scrape workers busy longer on slow targets.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/values.yaml`
- `docs/public/installation_guide.md`


### Exporter cannot authenticate to MongoDB

**Symptoms:**

- MongoDB exporter is installed and scraped, but MongoDB series stay empty.
- The monitoring user is missing or lacks the documented role.

**Root cause:**

The exporter is scraped, so this is failed query execution, not a missing ServiceMonitor. Docs name
`prometheusExporter.monitoringUser` and `monitoringUserRole` (default user `monitoring`, role `readWrite` on `test`
plus `clusterMonitor`). Do not read Secret values during diagnosis.

**Applies to:** mongodb-operator when the MongoDB exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `mongodb-prometheus-exporter` is installed and scraped.
2. Request non-secret evidence that the monitoring user exists and holds the documented role. Do not read Secret
   values.

If those artifacts are missing, request them as Data required.

**How to fix:**

1. Confirm without reading Secrets that the monitoring user exists and holds the documented role.
2. Reconcile mongodb-services so that user is created with `prometheusExporter.monitoringUserRole` when that key is
   present.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Granting clusterMonitor and readWrite widens that user's catalog access. Do not copy Secret values into the
diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/mongo_monitoring_user.yaml`
- `docs/public/configuration.md`
- `docs/public/troubleshooting.md`
- `docs/public/installation_guide.md`


### MongoDB ServiceMonitor basicAuth mismatch

**Symptoms:**

- ServiceMonitor `mongodb-prometheus-exporter-service-monitor` exists.
- Scrape of the exporter fails authentication.

**Root cause:**

The ServiceMonitor `basicAuth` references Secret `mongodb-prom-exporter-credentials.v1` (values
`prometheusExporter.exporterUser`, `exporterPassword`, `prometheusExporterSecretName`). A missing or mismatched
exporter credential is a scrape-auth failure, not a namespace selector miss.

**Applies to:** mongodb-operator when the MongoDB exporter ServiceMonitor exists. Image and chart versions are
supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `mongodb-prometheus-exporter-service-monitor` exists and scrape fails with an auth error.
2. Name exporter credential keys only when non-secret evidence shows `prometheusExporter.exporterUser` or the Secret
   name, never Secret values.

**How to fix:**

1. Align `prometheusExporter.exporterUser` and the ServiceMonitor basicAuth Secret with the exporter credentials
   without copying Secret values into the report.
2. Do not change Monitoring Operator namespace selectors for an auth failure on an existing monitor.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Changing exporter basicAuth credentials can briefly fail scrapes until both the exporter and the
ServiceMonitor use the same Secret.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/templates/secrets/prometheus-exporter-secret.yaml`
