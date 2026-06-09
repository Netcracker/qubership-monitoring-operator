# CoreDNS

A dashboard for the CoreDNS DNS server.

## Tags

* `dns`
* `coredns`
* `k8s`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Requests (total) | Shows total requests by protocol (udp, tcp) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Failed config reload | Shows count of failed config CoreDNS reload by pod | Default:<br/>Mode: absolute<br/>Level 1: 2<br/><br/> |  |
| Version | Shows CoreDNS version | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Enabled plugins | Shows enabled plugins for CoreDNS |  |  |
| Responses (duration) | Shows request duration | Default:<br/>Mode: absolute<br/>Level 1: 0.1<br/><br/> |  |
| CoreDNS panics | Shows total number of panics. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Requests (by qtype) | Shows count of DNS request by type | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Requests (by zone) | Shows count of requests by zone |  |  |
| Responses (by rcode) | Count of responses grouped by code | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Responses (size, udp/tcp) | Response size for "udp/tcp" protocol | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Requests (size, udp/tcp) | Request size for "udp/tcp" protocol | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Cache (hitrate) | Counter of cache hits by cache type. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Cache (count) | Total elements in the cache by cache type. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
