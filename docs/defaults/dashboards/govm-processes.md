# Go Processes

Process status published by Go Prometheus client library, e.g. memory used, fds open, GC details

## Tags

* `process`
* `golang`

## Panels

### Help

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Go Processes Help | Show a short help which describe dashboard and any specific cases show data |  |  |
<!-- markdownlint-enable line-length -->

### Quick facts

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Uptime | Shows time from VM start | Default:<br/>Mode: absolute<br/>Level 1: 60<br/>Level 2: 180<br/><br/> |  |
| Start time | Shows VM start time |  |  |
| Memory | Shows resident memory |  |  |
| Threads | Shows total threads count in selected area |  |  |
| Goroutines | Shows total routines count in selected area |  |  |
| Open FDs | Shows amount of currently opened file descriptors |  |  |
| Go Version | Show Go version for selected pod and container |  |  |
<!-- markdownlint-enable line-length -->

### Memory and GC

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Process memory | Shows resident and virtual memory volume | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Process memory / second | Shows per-second memory rate | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Go memstats | Shows how much Go is using from the amount of memory the Prometheus process is using from the kernel | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Go memstats / second | Shows per-second rate of how much Go is using from the amount of memory the Prometheus process is using from the kernel | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Goroutines | Shows the number of goroutined that currently exists | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| GC duration quantiles | Shows  the summary of the GC invocation durations | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Open fds | Shows the current number of used file descriptors | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Open / closed fds delta | Shows delta of used file descriptors. A positive value shows newly opened descriptors, a negative value shows closed descriptors. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### HTTP Statistics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| HTTP request by code: 2xx or 3xx $pod/$container | Shows amount of successful requests to different URLs per second | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| 95%-tile latency of response time | How long it took to process a request, partitioned by method and HTTP path | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### HTTP Errors

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| HTTP requests by code: 4xx $pod/$container | Shows amount of requests to different URLs per second with code 4xx | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| HTTP requests by code: 5xx $pod/$container | Shows amount of requests to different URLs per second with code 5xx | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### GRPC Statistics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Unary requests by OK code | Shows amount of successful requests to different methods per second | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Percentage unary requests by OK code | Shows percentage of successful requests to different methods per second | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| 95%-tile latency of unary response time | How long it took to process an unary request, partitioned by method | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### GRPC Errors

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Unary failed requests | Shows amount of error requests to different methods per second partitioned by error codes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Percentage unary failed requests | Shows percentage of failed requests to different methods per second | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### SQL Statistics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Open Connections | Show count of open SQL connections (statuses: idle, in use) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Blocked Connections | Show blocked connections and duration of blocking | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Closed connections | Shows number of closed SQL connections by reason | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
