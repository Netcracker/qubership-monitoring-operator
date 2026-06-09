# Request Handling Performance

## Tags

* `nginx`
* `ingress`
* `k8s`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total request handling time | Total time taken for nginx and upstream servers to process a request and send a response |  |  |
| Upstream response time | The time spent on receiving the response from the upstream server |  |  |
| Request volume by Path | Shows requests per second, group by path |  |  |
| Median upstream response time by Path | For each path observed, its median upstream response time |  |  |
| Response error rate by Path | Percentage of 4xx and 5xx responses among all responses. |  |  |
| Upstream time consumed by Path | For each path observed, the sum of upstream request time |  |  |
| Response error volume by Path | Show request errors with 4xx, 5xx codes by path |  |  |
| Average response size by Path | Show average response size by path |  |  |
| Upstream service latency | Show average upstream service latency |  |  |
<!-- markdownlint-enable line-length -->
