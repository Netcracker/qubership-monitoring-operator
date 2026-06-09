# Kubernetes / Pods Distribution By Zones

This dashboard shows extended statistics about distribution pods by zone

## Tags

* `k8s`
* `pods`
* `zones`

## Panels

### Zones overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of nodes $zone | Shows number of nodes for $zone |  | **Panel is multiplied by parameter `zone`** |
| Number of CPUs $zone | Shows number of  CPU for $zone |  | **Panel is multiplied by parameter `zone`** |
| Total RAM $zone | Shows amount of memory for $zone |  | **Panel is multiplied by parameter `zone`** |
| Used RAM $zone | Shows  amount of used memory for $zone |  | **Panel is multiplied by parameter `zone`** |
| Total number of pods $zone | Shows number of pods for $zone |  | **Panel is multiplied by parameter `zone`** |
<!-- markdownlint-enable line-length -->

### Pods distribution

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Ready Deployment Pods distribution per zones | Shows number of pods with ready=true status per Zones for kind: Deployment |  |  |
| Ready StatefulSet Pods distribution per zones | Shows number of pods with ready=true status per Zones for Kind: StatefulSet |  |  |
| Ready DaemonSet Pods distribution per zones | Shows number of pods with ready=true status per Zones for Kind: DaemonSet |  |  |
<!-- markdownlint-enable line-length -->
