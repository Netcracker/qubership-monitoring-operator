# Infrastructure component troubleshooting


## DRNavigator

### Site-manager monitor is not rendered

**Symptoms:**

- SM Monitoring dashboard is empty.
- Deployment `site-manager` exists, but ServiceMonitor `site-manager-service-monitor` and/or GrafanaDashboard CR
  `site-manager-grafana-dashboard` are absent.

**Root cause:**

The site-manager workload always binds metrics on port `metrics` (9000). `MONITORING_ENABLED` (default true) only
creates the ServiceMonitor and GrafanaDashboard CRs. This is missing scrape and dashboard resources, not a missing
site-manager Deployment.

**Applies to:** DRNavigator when site-manager Prometheus metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `site-manager` exists and `site-manager-service-monitor` is absent.
2. Name `MONITORING_ENABLED` only when effective or rendered values show it is false.

**How to fix:**

1. Set `MONITORING_ENABLED: true` when that key is present, then upgrade.
2. Confirm that `site-manager-service-monitor` selects `app: site-manager`, port `metrics`, and path `/metrics`.

Do not change Monitoring Operator selectors until the ServiceMonitor exists.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling monitoring CRs adds scrape of site-manager `/metrics` and loads the SM Monitoring dashboard.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/service-monitor.yaml`
- `charts/site-manager/templates/grafana-dashboard.yaml`
- `documentation/public/installation.md`


### paas-geo-monitor is not installed

**Symptoms:**

- Paas-Geo-Monitor Monitoring dashboard is empty, or series `peer_*` / `paas_geo_monitor_*` are absent.
- Deployment `paas-geo-monitor` is absent.

**Root cause:**

paas-geo-monitor is a second workload gated by `paasGeoMonitor.install` (default false). Stop here; do not treat an
absent geo ServiceMonitor as a scrape failure while the workload is missing.

**Applies to:** DRNavigator when paas-geo-monitor metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `paas-geo-monitor` is absent.
2. Name `paasGeoMonitor.install` only when effective values show it is false.

**How to fix:**

1. Set `paasGeoMonitor.install: true` when that key is present, then upgrade.
2. Confirm that Deployment and Service `paas-geo-monitor` exist, with port `web`.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling paas-geo-monitor deploys a connectivity probe workload.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/paas-geo-monitor-deployment.yaml`
- `documentation/public/installation.md`


### paas-geo-monitor scrape is not rendered

**Symptoms:**

- Deployment `paas-geo-monitor` exists.
- ServiceMonitor `paas-geo-monitor-service-monitor` and/or GrafanaDashboard CR `paas-geo-monitor-grafana-dashboard` are
  absent.

**Root cause:**

Geo ServiceMonitor and dashboard templates require both `paasGeoMonitor.install` and `MONITORING_ENABLED`. This is
missing scrape and dashboard resources for an installed geo workload.

**Applies to:** DRNavigator when paas-geo-monitor is installed. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `paas-geo-monitor` exists and `paas-geo-monitor-service-monitor` is absent.
2. Name `MONITORING_ENABLED` only when effective values show it is false.

**How to fix:**

1. Set `MONITORING_ENABLED: true` when that key is present, then upgrade.
2. Confirm that `paas-geo-monitor-service-monitor` selects `app: paas-geo-monitor`, port `web`, and path `/metrics`.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling geo monitoring CRs adds a scrape target for paas-geo-monitor.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/templates/paas-geo-monitor-service-monitor.yaml`
- `charts/site-manager/templates/grafana-dashboard.yaml`


### paas-geo-monitor peer series are absent

**Symptoms:**

- paas-geo-monitor ServiceMonitor exists and is scraped.
- Series `peer_dns_status`, `peer_svc_status`, or `peer_pod_status` are absent, while health or CPU panels may work.

**Root cause:**

Labeled peer gauges are written only by the ping loop over configured peers. Docs mark
`paasGeoMonitor.env.PAAS_PING_PEERS` as required for monitoring (default true). `paasGeoMonitor.config.peers` defaults
to `[]`, so no peer series are emitted until peers are configured.

**Applies to:** DRNavigator when paas-geo-monitor is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `paas-geo-monitor-service-monitor` is scraped.
2. Name `PAAS_PING_PEERS` or `paasGeoMonitor.config.peers` only when effective values show ping is off or the peer list
   is empty.

**How to fix:**

1. Set `paasGeoMonitor.env.PAAS_PING_PEERS: true` and configure `paasGeoMonitor.config.peers` when those keys are
   present, then upgrade.
2. Do not set `paasGeoMonitor.install: true` as the fix for missing peer series on a healthy scrape.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling peer pings sends DNS, Service, and pod probes to each configured peer.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `paas-geo-monitor/pkg/app/app.go`
- `documentation/public/installation.md`


### paas-geo-monitor BGP series are absent

**Symptoms:**

- paas-geo-monitor target is healthy.
- Series `paas_geo_monitor_bgp_peer` or `paas_geo_monitor_bgp_route` are absent, and BGP or Service Network panels stay
  empty.

**Root cause:**

BGP series register only when env `PAAS_BGP_METRICS` is `true` (`paasGeoMonitor.env.PAAS_BGP_METRICS`, values default
`"true"`). This is a disabled series set, not missing peer ping.

**Applies to:** DRNavigator when paas-geo-monitor is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that the geo target is healthy and BGP series are absent.
2. Name `PAAS_BGP_METRICS` only when effective values or env show it is not `true`.

**How to fix:**

1. Set `paasGeoMonitor.env.PAAS_BGP_METRICS: "true"` when that key is present, then upgrade.
2. Do not treat this as an empty `paasGeoMonitor.config.peers` list.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling BGP metrics queries BGP state on the node.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/paas-geo-monitor-deployment.yaml`
- `paas-geo-monitor/pkg/app/app.go`


### paas-geo-monitor service CIDR dashboard regex mismatch

**Symptoms:**

- BGP series exist.
- Service Network Availability variable or panels stay empty because `$service_cidr` does not match.

**Root cause:**

Dashboard variable `service_cidr` uses regex from `paasGeoMonitor.monitoring.serviceCIDRRegex` (default `172.*`). A
cluster service subnet outside that regex leaves those panels empty. This is a dashboard query mismatch, not missing
BGP series.

**Applies to:** DRNavigator when paas-geo-monitor BGP series exist. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that BGP series exist and the empty panels use `$service_cidr`.
2. Name `paasGeoMonitor.monitoring.serviceCIDRRegex` only when effective values or the dashboard JSON show a regex that
   does not match the service subnet.

**How to fix:**

1. Set `paasGeoMonitor.monitoring.serviceCIDRRegex` to the service subnet regex when that key is present, then
   upgrade.
2. Do not enable `PAAS_BGP_METRICS` as the fix for a regex mismatch.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** A broader CIDR regex shows more BGP service networks in Grafana.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/grafana-dashboard.yaml`
- `charts/site-manager/monitoring/paas-geo-monitor-grafana-dashboard.json`
