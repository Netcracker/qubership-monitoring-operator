# Prometheus Self Monitoring

There is no description on this dashboard

## Tags

* `k8s`
* `prometheus`
* `self-monitor`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Uptime | Show the prometheus instance uptime | Default:<br/>Mode: absolute<br/>Level 1: 3600<br/><br/> |  |
| Last Successful Config Reload | Show last successful config reload |  |  |
| CPU usage | Shows CPU usage for the prometheus instance | Default:<br/>Mode: absolute<br/>Level 1: 0.85<br/>Level 2: 1.2<br/>Level 3: 2<br/><br/> |  |
| Resident memory | The amount of memory the Prometheus process is using from the kernel | Default:<br/>Mode: absolute<br/>Level 1: 6442450944<br/>Level 2: 9663676416<br/><br/> |  |
| Data | The number of bytes that are currently used for local storage by all blocks |  |  |
| Oldest data | Show the earliest saved data |  |  |
| Count of Jobs | Shows the number of active jobs | Default:<br/>Mode: absolute<br/>Level 1: 300<br/>Level 2: 500<br/><br/> |  |
| Series | Shows the number of active time series | Default:<br/>Mode: absolute<br/>Level 1: 500000<br/>Level 2: 1000000<br/><br/> |  |
| Points per second | Show number of metrics per second stored by Prometheus as taken from the last $interval. |  |  |
| Build, uptime and runtime instance info | Prometheus instance overview | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Main info

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of time series | Shows the number of active time series |  |  |
| $quantile quantile of interval length between scrapes per job | Actual amount of time between target scrapes |  |  |
| Prometheus errors in $interval | Shows the total number of Prometheus errors, including different non-default behavior. <br/><br/>For example, if you are sending not all the metrics to the remoteStorage, then the panel will also display the number of unsent points (prometheus_remote_storage_samples_dropped_total). |  |  |
<!-- markdownlint-enable line-length -->

### Exporters/Targets

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Ready jobs/ All jobs | Shows the number of jobs(Ready/All) |  |  |
| kubelet | Show number of targets for kubelet | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| node-exporter | Show number of targets for node-exporter | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| kube-state-metrics | Show number of targets for kube-state-metrics | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| etcd | Show number of targets for etcd | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| kube-apiserver | Show number of targets for kube-apiserver | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| nginx-ingress | Show number of targets for nginx-ingress<br/><br/>* Works only in kubernetes | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Scrape targets | Show current number of targets in this scrape pool |  |  |
| Discovered targets | Show current number of discovered targets |  |  |
<!-- markdownlint-enable line-length -->

### Cardinality

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of series by job | Show number of series by job |  |  |
| New series by job | Show how many new series were created |  |  |
| Highest cardinality for $job | Show 10 metrics with the highest cardinality for $job target | Default:<br/>Mode: absolute<br/>Level 1: 1000<br/>Level 2: 5000<br/>Level 3: 10000<br/><br/> | **Panel is multiplied by parameter `job`** |
| Count metrics for $job | Show number of metrics for $job | Default:<br/>Mode: absolute<br/>Level 1: 5000<br/>Level 2: 10000<br/>Level 3: 50000<br/><br/> | **Panel is multiplied by parameter `job`** |
<!-- markdownlint-enable line-length -->

### Requests & queries

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| HTTP requests/s | Average request duration for the last 1m. |  |  |
| HTTP request latency | Latencies for HTTP requests |  |  |
| Time spent in HTTP requests/s | Time spent in HTTP requests/s |  |  |
| HTTP request count by handler in $interval | Counter of HTTP requests |  |  |
| $quantile quantile of HTTP request duration per handler | HTTP request duration per handler |  |  |
| $quantile of request size by handler | Request size by handler |  |  |
<!-- markdownlint-enable line-length -->

### Resources

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU usage | Shows CPU usage for the prometheus instance |  |  |
| Memory usage | Resident and allocated memory size in bytes.<br/> |  |  |
| Allocations per second in $interval | Total number of bytes allocated, even if freed |  |  |
<!-- markdownlint-enable line-length -->

### TSDB stats

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Oldest data | Show the earliest saved data |  |  |
| Size of the storage | The number of bytes that are currently used for local storage by all blocks |  |  |
<!-- markdownlint-enable line-length -->

### Query engine & API

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of concurent queries and it's limit | The current and max number of queries being executed or waiting |  |  |
| $quantile quantile of query engine evaluation duration per slice | Query engine evaluation duration per slice |  |  |
| Number of queries on remote read API | The current number of remote read queries being executed or waiting |  |  |
<!-- markdownlint-enable line-length -->

### Rule evaluation

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Percentage rule group evaluation duration from the rule group evaluation interval (top 20) | Percentage rule group evaluation duration from the rule group evaluation interval (top 20) |  |  |
| $quantile quantile of rule evaluation duration | Rule evaluation duration (for quantile => 0.5) |  |  |
| Number of rules per group (top 20 groups) | Number of rules for the first 20 groups |  |  |
<!-- markdownlint-enable line-length -->

### Alerting

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Notification queue size and capacity | The number of alert notifications in the queue. And the capacity of the alert notifications queue. |  |  |
| Number of sent notifictions per alertmanager in $interval | Total number of alerts sent |  |  |
| $quantile of notification latency per alertmanager | Latency quantiles for sending alert notifications |  |  |
<!-- markdownlint-enable line-length -->

### Scraping

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Scrape Interval | Intervals between scrapes |  |  |
| Number of target errors | Number of scrapes that hit the sample limit and were rejected. And the number of samples rejected due to duplicate timestamps but different values, not being out of the expected order, timestamp falling outside of the time bounds.  |  |  |
| Metadata Cache Size | The number of bytes that are currently used for storing metric metadata in the cache |  |  |
<!-- markdownlint-enable line-length -->

### Service discovery

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of discovered targets per Type and config | Current number of discovered targets |  |  |
| Service discovery  sync count by mechanism in $interval | Service discovery  sync count by mechanism |  |  |
| $quantile quantile of refresh duration per SD mechanism | <br/> |  |  |
<!-- markdownlint-enable line-length -->

### Compaction and retention

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| $quantile quantile of compaction duration | Duration of compaction runs |  |  |
| Number of vertical/horizontal compactions in last $interval | Number of vertical/horizontal compactions |  |  |
| Number of time/size retention cutoffs in $interval | The number of times that blocks were deleted because the maximum time limit or number of bytes was exceeded. |  |  |
<!-- markdownlint-enable line-length -->

### WAL

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| $quantile quantile of WAL fsync duration | Duration of WAL fsync ( for quantile => 0.5 ) |  |  |
| Number of completed pages in $interval | Number of completed pages |  |  |
| Average duration of WAL truncation | *TODO: Fill panel description* |  |  |
<!-- markdownlint-enable line-length -->

### Remote storage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of queries to remote storage per client | Number of remote read queries |  |  |
| Number of samles successfuly sent to remote storage per queue | Total number of samples successfully sent to remote storage. |  |  |
| $quantile quantile of batch send call duration to remote storage | Call duration to remote storage |  |  |
| Number of used shards for concurrent sending of data to remote storage and it's capacity | Number of used shards for concurrent sending of data to remote storage and it's capacity |  |  |
<!-- markdownlint-enable line-length -->

### Go stats

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of gorutines | Number of goroutines that currently exist. |  |  |
| Duration of Go garbage collection | Max duration of garbage collection cycles. |  |  |
| Go system memory allocations | Go system memory allocations |  |  |
<!-- markdownlint-enable line-length -->

### Network

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Accepted/closed inbound connections per listener | Accepted/closed inbound connections per listener of a given name. |  |  |
| Established/closed outbound connections per dialer  | Established/closed outbound connections per dialer a given name. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Rule Group: $RuleGroup

**Row Rule Group: $RuleGroup is multiplied by parameter `RuleGroup`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| $RuleGroup: Duration | The interval and last duration evaluation of a rule group.<br/> |  |  |
| $RuleGroup: Rules | The number of rules |  |  |
<!-- markdownlint-enable line-length -->
