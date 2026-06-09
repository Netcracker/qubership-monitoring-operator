# Kubernetes / API Server

There is no description on this dashboard

## Tags

* `k8s`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| API server status | Shows API server status depending on the number of endpoints exposing metrics | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 2<br/><br/> |  |
| API Servers | Shows the number of running API servers  depending on the number of endpoints exposing metrics | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 2<br/><br/> |  |
| Current Error Rate | Show the current API server error rate: percentage of requests with errors out of the total number of requests | Default:<br/>Mode: absolute<br/>Level 1: 3<br/><br/> |  |
| Average API Server Latency | Shows the average duration of request to API server | Default:<br/>Mode: absolute<br/>Level 1: 0.2<br/><br/> |  |
| Average Etcd Requests Latency | Shows the average duration of requests to Etcd | Default:<br/>Mode: absolute<br/>Level 1: 0.2<br/><br/> |  |
| Current Controller work queue Depth | Show the current depth of controller work queue | Default:<br/>Mode: absolute<br/>Level 1: 5<br/><br/> |  |
| API server nodes status | Shows API server nodes status |  |  |
| Request Rate | Shows the number of requests per second by each API server instance |  |  |
<!-- markdownlint-enable line-length -->

### Requests

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Request rate by verbs | Shows the number of requests per second by verb |  |  |
| Request Rate by resource and group | Shows the number of requests per second by resource and its API group in format: API/resource |  |  |
<!-- markdownlint-enable line-length -->

### Latency

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Request duration 0.99 quantile | Shows the request duration  0.99 quantile: 99% of requests has equal or less duration  |  |  |
| Request duration by verb 0.99  quantile | Shows the request duration  0.99 quantile by verb: 99% of requests has equal or less duration  |  |  |
| Request duration by resource and group 0.99 quantile | Shows the request duration  0.99 quantile by resource type and its API group: 99% of requests  . Format: API/resource<br/> |  |  |
<!-- markdownlint-enable line-length -->

### Errors

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| API Server error percentage | Shows the percentage of requests failed with 5xx error by each API server instance |  |  |
| Error request rate by verb | Show the number of requests per second failed with 5xx or 4xx code grouped by verb |  |  |
| Error request rate by code | Show the number of requests per second failed with 5xx or 4xx code grouped by code |  |  |
<!-- markdownlint-enable line-length -->

### Controller Work Queue

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Adds to work queue per second | Shows the number of adds to controller work queue per second |  |  |
| Work queue depth | Shows the current depth of controller work queue. It should be near 0 |  |  |
| Item processing time 0.99  quantile | Shows 0.99 duration quantile of the item processed from work queue: 99% of items has equal or less processing time |  |  |
| Item queue time 0.99 quantile | Shows 0.99 duration quantile of the item stay work queue: 99% of items has equal or less staying time |  |  |
| Unfinished work time | Show time of unfinished item processing.  Large values indicate stuck threads |  |  |
| Longest running processor time | Show longest time of item processing |  |  |
<!-- markdownlint-enable line-length -->

### Etcd

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Etcd request rate | Shows the number of requests per second sent to Etcd |  |  |
| Etcd request duration 0.99 quantile | Shows the Etcd request duration  0.99 quantile: 99% of requests has equal or less duration  |  |  |
| Request duration 0.99 quantile by operation | Shows the Etcd request duration  0.99 quantile by operation: 99% of requests has equal or less duration  |  |  |
| Request duration 0.99 quantile by object type | Shows the Etcd request duration  0.99 quantile by object type: 99% of requests has equal or less duration  |  |  |
| Objects stored in etcd | Shows the number of objects stored in Etcd |  |  |
<!-- markdownlint-enable line-length -->

### Go stats

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU usage | Shows the CPU usage by each API server instance |  |  |
| Goroutines | Shows the number of goroutines by each API server instance |  |  |
| Resident memory | Show the resident memory usage by each API server instance |  |  |
| Allocated memory | Shows the allocated  memory  by each API server instance |  |  |
| Allocations per second | Shows the memory allocations per second by each instance |  |  |
<!-- markdownlint-enable line-length -->
