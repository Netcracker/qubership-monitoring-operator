# TLS Status

Shows TLS statuses for all services in the cluster

## Tags

* `k8s`
* `tls`

## Panels

### Help

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Help | A detailed description for the dashboard |  |  |
<!-- markdownlint-enable line-length -->

### TLS Status

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| All services | Number of all services<br/><br/>(pods with unique `app.kubernetes.io/name` labels) |  |  |
| List of TLS Statuses | Shows names of services (by `app.kubernetes.io/name` label) and their TLS statuses (by `app.kubernetes.io/tls` label) |  |  |
| TLS enabled | Number of services with enabled TLS<br/><br/>(pods have `app.kubernetes.io/tls: true` label) | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| TLS disabled | Number of services with disabled TLS<br/><br/>(pods have `app.kubernetes.io/tls: false` label) | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| TLS status unknown | Number of services with unknown TLS status<br/><br/>(pods with empty or null `app.kubernetes.io/tls` label) | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->
