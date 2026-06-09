# Kafka Java Clients Monitoring

Kafka Java Producers and Consumers Monitoring Dashboard

## Tags

* `Prometheus`
* `Kafka`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Producers Overview | Number of records and bytes produced per second for a specific topic or across all topics. |  |  |
| Consumers Overview | Number of records and bytes consumed per second for a specific topic or across all topics. |  |  |
<!-- markdownlint-enable line-length -->

### Resources Metrics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Usage | CPU usage by pod. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Memory Usage | Memory usage by pod. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Heap Memory Usage | Heap memory usage by pod. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Network Rate | Network Rate: transmitted (tx) and received (rx). | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Producers Metrics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Producer Records Rate | The rate of records that successfully processed by Producer and written to Kafka topics. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Records Bytes Rate | The rate of bytes that successfully processed by Producer and written to Kafka topics. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Request Rate | Average number of requests sent per second. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Errors Count | The number of errors during producing. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Producer Record Queue Time | The average time in milliseconds record batches spent in the send buffer. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Producer Request Latency | Average request latency (in ms). | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Errors Rate | The rate of errors during producing. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Records Size | The average record size. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Response Count | Number of responses received. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Topics Errors Count | The number of errors during producing by topics. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Records Batch Size | The average record batch size. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Response Rate | Average number of responses received per second. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Waiting Threads | The number of threads in waiting state. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average IO Wait Producer Time | Average length of time the I/O thread spent waiting for a socket (in ns). | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Compression Rate | Average compression rate of batches sent. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| In-Flight Requests | The number of in-flight requests. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Producer Outgoing Bytes Rate | The average number of outgoing/incoming bytes per second. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Consumer Metrics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Consumer Records Rate | Average number of records consumed per second for a specific topic or across all topics. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumer Records Bytes Rate | Average number of bytes consumed per second for a specific topic or across all topics. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumer Request Rate | The rate of requests sent. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumer Lag | Number of messages consumer is behind producer on this partition. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Fetch Rate | The number of fetch requests per second from the consumer. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Consumer Request Size | The average Consumer request size. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumer Rebalance Count | Number of rebalances and failed rebalances. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| IO Wait Consumer Time Rate | Average length of time the I/O thread spent waiting for a socket (in ns). | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Consumer Records Per Request | The average number of records per request. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumer Rebalance Rate | The rate of rebalances and failed rebalances. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumer Incoming Bytes Rate | The number of bytes that successfully processed by Consumer and written to Kafka topics per second. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average Time Between Poll | The average time between poll in ms. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Listener Metrics (Spring)

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Listener Calls Rate | Shows per-second rate of listener calls count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Listener Failures Rate | Shows the count of Kafka listener failures |  |  |
| Average Listener Processing Duration | Shows the average time of successful call processing by listener | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
