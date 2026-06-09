# Prometheus / Graphite remote adapter

There is no description on this dashboard

## Tags

* `k8s`
* `graphite`

## Panels

### Processing

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Received / Sent / Failed (Instant) | Total number of received samples.<br/>Total number of processed samples sent to remote <br/>storage.<br/>Total number of processed samples which failed on send to remote storage. |  |  |
| Received / Sent / Failed | Total number of received samples.<br/>Total number of processed samples sent to remote <br/>storage.<br/>Total number of processed samples which failed on send to remote storage.<br/> |  |  |
| Received Samples | Total number of received samples |  |  |
| Processed Samples | Total number of processed samples sent to remote storage |  |  |
| Failed samples | Total number of processed samples which failed on send to remote storage |  |  |
| Sent samples duration | Duration of sample batch send calls to the remote storage |  |  |
<!-- markdownlint-enable line-length -->

### Go process

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Usage | Total user and system CPU time spent in seconds. |  |  |
| Process memory | Shows resident and virtual memory volume |  |  |
| Process memory / second | Shows per-second memory rate |  |  |
| Go memstats | Shows how much Go is using from the amount of memory the Prometheus process is using from the kernel |  |  |
| Go memstats / second | Shows per-second rate of how much Go is using from the amount of memory the Prometheus process is using from the kernel |  |  |
| Goroutines | Shows the number of goroutined that currently exists |  |  |
| GC duration quantiles | Shows  the summary of the GC invocation durations |  |  |
<!-- markdownlint-enable line-length -->

### API Requests

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Handle API request duration, 99th quantile | A histogram of latencies for requests. |  |  |
| API Response size, 99th quantile | A histogram of response sizes for requests. |  |  |
<!-- markdownlint-enable line-length -->

### Sent metrics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Sent batch duration | Duration of sample batch send calls to the remote storage. |  |  |
| Sent samples | Total number of processed samples sent to remote storage. |  |  |
<!-- markdownlint-enable line-length -->

### Promhttp

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Handled requests | Total number of scrapes by HTTP status code. |  |  |
| Requests in-flight | Current number of scrapes being served. |  |  |
<!-- markdownlint-enable line-length -->
