# Version overview

There is no description on this dashboard

## Tags

* `version-exporter`
* `k8s`

## Panels

### Kubernetes

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| kubernetes | Shows kubernetes build and node info | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| etcd  | Shows server version info | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| coredns | Shows build info. In kubernetes only. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| nginx ingress | Shows build info. In kubernetes only. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Monitoring services

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Grafana | Shows build info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Prometheus | Shows build info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Prometheus Operator | Shows build info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Alert Manager | Shows build info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Node Exporter | Shows build info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Version Exporter | Shows build info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Containers

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Container Image | Shows pod container info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Init Container | Shows build pod init container info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Runtime

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Golang | Shows go info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Python | Shows python info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| JVM | Shows jvm info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Version exporter metrics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Info | *TODO: Fill panel description* |  |  |
| $version_metrics | Shows collected $version_metrics metric info. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `version_metrics`** |
<!-- markdownlint-enable line-length -->
