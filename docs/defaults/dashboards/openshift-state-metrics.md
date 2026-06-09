# Openshift / State metrics

Dashboards shows openshift-state-metrics specific metrics

## Tags

* `openshift-state-metrics`
* `openshift`

## Panels

### Deploymentconfigs

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Deploymentconfigs | Shows the counts of running/not running deploymentconfigs in the cluster |  |  |
| Deploymentconfigs list | Shows a full list of Openshift deploymentconfigs in the cluster | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Resource Quotas Requests

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Percent CPU used vs request quota | Shows  CPU used vs request quota in percents | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 2<br/><br/> |  |
| Percent Memory Used vs Request Quota | Shows memory used vs Request quota in percent | Default:<br/>Mode: absolute<br/>Level 1: 0.2<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Resource Quotas Limits

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Percent CPU Used vs Limit Quota | Shows CPU Used vs Limit Quota in percent | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Percent Memory Used vs Limit Quota | Shows t Memory Used vs Limit Quota in percent | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Routes

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| List of routes | Shows list of all the routes present in the cluster, aggregated by host and namespace | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Duplicated routes | Shows duplicated routes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Routes with errors | Shows openshift routes with errors | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Builds

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| New Builds | Shows New Builds by Processing Time | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Builds with Errors | Show openshift builds with error | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->