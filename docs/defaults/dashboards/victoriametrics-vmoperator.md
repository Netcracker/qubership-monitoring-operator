# Victoria Metrics / Operator

Overview for operator VictoriaMetrics v0.25.0 or higher

## Tags

* `self-monitor`
* `victoriametrics`
* `vmoperator`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Version | Victoria Metrics operator version. |  |  |
| CRD Objects count by controller | Number of objects at kubernetes cluster per each controller | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Uptime | Victoria Metrics operator uptime. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Reconciliation rate by controller | Total number of reconciliations per controller. |  |  |
| Log message rate | Victoria Metrics operator log message rate. |  |  |
<!-- markdownlint-enable line-length -->

### Troubleshooting

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| reconcile errors by controller | Non zero metrics indicates about error with CR object definition (typos or incorrect values) or errors with kubernetes API connection. |  |  |
| throttled reconcilation events | Operator limits number of reconcilation events to 5 events per 2 seconds.<br/> For now, this limit is applied only for vmalert and vmagent controllers.<br/> It should reduce load at kubernetes cluster and increase operator performance. |  |  |
| Wokring queue depth | Number of objects waiting in the queue for reconciliation. Non-zero values indicate that operator cannot process CR objects changes with the given resources. |  |  |
| Reconcilation latency by controller |  For controllers with StatefulSet it's ok to see latency greater then 3 seconds. It could be vmalertmanager,vmcluster or vmagent in statefulMode.<br/><br/> For other controllers, latency greater then 1 second may indicate issues with kubernetes cluster or operator's performance.<br/>  |  |  |
<!-- markdownlint-enable line-length -->

### Resources

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Memory usage ($instance) | Victoria Metrics operator memory usage. |  |  |
| CPU ($instance) | Victoria Metrics operator CPU usage. |  |  |
| Goroutines ($instance) | Total number of goroutines. |  |  |
| GC duration ($instance) | Victoria Metrics operator avg GC duration. |  |  |
<!-- markdownlint-enable line-length -->
