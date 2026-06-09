# Kubernetes / Distribution by Labels

The dashboard allows filtering Kubernetes resources (pods, ingresses, etc.) by app.kubernetes.io labels

## Tags

* `k8s`
* `labels`

## Panels

### Help

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Help | Shows information about current dashboard |  |  |
<!-- markdownlint-enable line-length -->

### Pods

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Pods | Shows pods with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Ingresses

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Ingresses | Shows ingresses with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Deployments

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Deployments and ReplicaSets | Shows deployments and replicasets with associated app.kubernetes.io labels, if they have more than 0 ready replicas | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### StatefulSets

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| StatefulSets | Shows statefulsets with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### DaemonSets

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| DaemonSets | Shows daemonsets with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Jobs

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Jobs | Shows Kubernetes jobs with associated app.kubernetes.io labels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
