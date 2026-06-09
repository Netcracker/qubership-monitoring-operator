# Openshift / HAProxy

Dashboard shows HAProxy specific metrics for Openshift-ingress namespace

## Tags

* `openshift`
* `openshift-haproxy`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Backends Up | Shows amount of instances for current backend with status "UP" | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Current health status of the backends | Shows states in time for current backend and route(if exists) | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Backends DOWN | Shows amount of instances for current backend with status "DOWN" | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Basic General Info

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total responses by HTTP code | Shows statistics about total responses from backends and frontends by HTTP codes |  |  |
| Current total of incoming / outgoing bytes | Shows traffic statistics for backends and frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Sessions | Shows amount of current sessions and rate for frontends and backends |  |  |
| Total number of connections | Sows amount of connections and errors(only to backends) to frontends and backends |  |  |
<!-- markdownlint-enable line-length -->

### Throughput

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Back - Current total of incoming / outgoing bytes | Shows detailed traffic statistics for backends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Front - Current total of incoming / outgoing bytes | Shows traffic statistics per frontend | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Connections

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Backends - Total number of connections OK / Error | Shows statistics about connections and errors per backend | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Frontends - Total number of connections | Shows amount of connections per frontend | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Queue

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Backends - Maximum observed number of queued requests  | Shows max amount of queued requests per route | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Backends - Current number of queued requests | Shows current amount of queued requests per route | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Requests / Responses

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Backends - Total of re-dispatch / error / retry warnings | Shows amount of problems(errors, retry warnings, re dispatch) per route | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Frontends- Total HTTP requests OK / Error / Denied | Shows amount of requests per status for each frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Responses by HTTP code

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Backends - Total of HTTP responses | Shows amount of responses for backends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Frontends - Total of HTTP responses | Shows amount of responses for frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Sessions Backends

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Backends - Current number of active sessions | Shows amount of current sessions for backends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Back - Current number of sessions rate per second over last elapsed second | Shows amount of sessions rate per second over last elapsed second for backends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Backends - Maximum observed number of active sessions and limit | Shows max  amount of active sessions and limit for backends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Backends - Maximum number of sessions rate per second | Shows max amount of sessions rate per second for backends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Sessions Frontends

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Front - Current number of active sessions | Shows amount  of active sessions for frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Front - Current number of sessions rate per second over last elapsed second | Shows amount of sessions rate per second over last elapsed second for frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Front - Maximum observed number of active sessions and limit | Shows max amount of active sessions and limit for frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Front - Maximum number of sessions rate per second and limit | Shows max amount of sessions rate per second and limit for frontends | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->