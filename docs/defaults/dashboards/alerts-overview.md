# Alerts - Overview

Show current firing and pending alerts, and  severity alert counts.

## Tags

* `alerts`
* `k8s`

## Panels

### Alerts Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Alerts | Show the number of alerts  (all alerts, sent to alertmanager, pending sent to alertmanager) | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Critical | Show the number of alerts with "critical" severity | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| High | Show the number of alerts with "high" severity | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Warning | Show the number of alerts with "warning" severity | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Info | Show the number of alerts with "info" severity | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Other Labels | Show the number of alerts without severity | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Alerts

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Firing | Show the table with alerts in firing state, which group by alertnames, namespaces, severities |  |  |
| Pending | Show the table with alerts in a pending state, which group by alertnames, namespaces, severities | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
