# Openshift / Cluster version operator

Dashboards shows Cluster version operator specific metrics

## Tags

* `openshift`
* `cluster-version-operator`

## Panels

### Cluster version overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Operator status | Shows cluster version operator status | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Current version | Shows current version of the cluster |  |  |
| Updated version | Shows version from which current was updated |  |  |
| Initial version of cluster | Shows initial version of the cluster |  |  |
| Payload | Shows statistics about applied, pending and errors |  |  |
<!-- markdownlint-enable line-length -->

### Cluster version capability

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Cluster capabilities | Shows currently enabled cluster capabilities | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Cluster operators

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Cluster operators status | Shows cluster operators status by  reason |  |  |
| Available  cluster operators | Available  cluster operators by reason |  |  |
| Degraded  cluster operators | Degraded cluster operators by reason |  |  |
| Upgradeable  cluster operators | Shows cluster operators which can be upgraded by reason |  |  |
<!-- markdownlint-enable line-length -->

### Go stats

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU usage | Shows the CPU usage by cluster version operator | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Goroutines | Shows the number of goroutines for cluster version operator | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Resident memory | Show the resident memory usage by cluster version operator | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Allocated memory | Shows the allocated  memory  by cluster version operator | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Allocations per second | Shows the memory allocations per second by  cluster version operator | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->