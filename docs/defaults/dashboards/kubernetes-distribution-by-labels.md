# Kubernetes / Distribution by Labels

The dashboard allows filtering Kubernetes resources (pods, ingresses, etc.) by app.kubernetes.io labels

## Tags

* `k8s`
* `labels`

## Panels

### Help

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Help | Shows information about current dashboard |  |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->

### Pods

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Pods | Shows pods with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->

### Ingresses

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Ingresses | Shows ingresses with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->

### Deployments

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Deployments and ReplicaSets | Shows deployments and replicasets with associated app.kubernetes.io labels, if they have more than 0 ready replicas | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->

### StatefulSets

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| StatefulSets | Shows statefulsets with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->

### DaemonSets

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| DaemonSets | Shows daemonsets with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->

### Jobs

<!-- markdownlint-disable line-length table-column-style no-space-in-emphasis -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Jobs | Shows Kubernetes jobs with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length table-column-style no-space-in-emphasis -->
