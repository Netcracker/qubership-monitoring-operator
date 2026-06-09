# Grafana Overview

Dashboards shows extended statistics about Grafana, including time of searching, saving and getting dashboards.

## Tags

* `self-monitor`
* `grafana`

## Panels

### Overview for $instance

**Row Overview for $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Dashboards count | Shows total amount of dashboards  | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Organisation count | Shows total amount of organisations | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Users count | Shows total amount of users | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Firing Alerts | Shows total amount of firing alerts | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Build Info | Shows statistics  about version and build of Grafana | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Responce Codes and Logins for $instance

**Row Responce Codes and Logins for $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total Response Statuses | Shoms total amount of http response status grouped by code and type | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Login Counts | Shows api login counters | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| RPS | Shows requests per second grouped by status code | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Request Latency | Shows requests latency include percentile  | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Dashboard Get Stats for $instance

**Row Dashboard Get Stats for $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Dashboard Get Duration Quantiles | summary for dashboard get duration | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Dashboard Get Duration Avg [5 min] | average dashboard get duration over 5 minutes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Dashboard Save Stats for $instance

**Row Dashboard Save Stats for $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Dashboard Save Duration Quantiles | Shows summary for dashboard save duration | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Dashboard Save Duration Avg [5 min] | Shows average dashboard save duration over 5 minutes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Dashboard Search Stats for $instance

**Row Dashboard Search Stats for $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Dashboard Search Duration Quantiles | Shows summary for dashboard search duration | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Dashboard Search Duration Avg [5 min] | Shows average dashboard search duration over 5 minutes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Datasource Stats for $instance

**Row Datasource Stats for $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total  Requests To Datasource | Show number of requests by code per datasource | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Datasource Request Duration Avg [5 min] | Average datasource request duration over 5 minutes by code per datasource | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Plugins on $instance

**Row Plugins on $instance is multiplied by parameter `instance`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Grafana plugins | Shows installed plugins | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
