# JVM Processes

Dashboard for Micrometer instrumented applications (Java, Spring Boot, Micronaut)

## Tags

* `process`
* `java`

## Panels

### Help

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| JVM Processes Help | *TODO: Fill panel description* |  |  |
<!-- markdownlint-enable line-length -->

### Quick Facts

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Uptime | Shows time from VM start | Default:<br/>Mode: absolute<br/>Level 1: 60<br/>Level 2: 180<br/><br/> |  |
| Start time | Shows VM start time |  |  |
| Heap used | Shows memory used by heap | Default:<br/>Mode: absolute<br/>Level 1: 70<br/>Level 2: 90<br/><br/> |  |
| Non-heap used | Shows memory used by non-heap | Default:<br/>Mode: absolute<br/>Level 1: 70<br/>Level 2: 90<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### JVM Memory

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| JVM Heap | Shows the heap JVM memory  used volume | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| JVM Non-Heap | Shows the non-heap JVM memory  used volume | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| JVM Total | Shows the amount of memory <br/> - used<br/> - committed to use<br/> -  max that can be used | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| JVM Process Memory | Shows process_memory_vss_bytes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### JVM Misc

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Usage | Shows system CPU usage,  process CPU usage and average process CPU usage | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Load | Shows<br/> - the sum of the number of runnable entities queued to available processors and the number of runnable entities running on the available processors averaged over a period of time<br/> - the number of processors available to the Java virtual machine<br/><br/> | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Threads | Shows<br/> - the current number of live threads including both daemon and non-daemon threads<br/> - the current number of live daemon threads<br/> - the peak live thread count since the Java virtual machine started or peak was reset<br/> -  the process threads | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Thread States | Shows the current number of threads having NEW state<br/> | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Log Events | Show increase of log events | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| File Descriptors | Show amount of processed files<br/> - open<br/> - max that can be open | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### JVM Memory Pools (Heap)

**Row JVM Memory Pools (Heap) is multiplied by parameter `persistence_counts`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| $jvm_memory_pool_heap | Shows amount of space<br/> - used<br/> - commited to be used<br/> - max available for use | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `jvm_memory_pool_heap`** |
<!-- markdownlint-enable line-length -->

### JVM Memory Pools (Non-Heap)

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| $jvm_memory_pool_nonheap | Shows amount of space<br/> - used<br/> - commited to be used<br/> - max available for use | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `jvm_memory_pool_nonheap`** |
<!-- markdownlint-enable line-length -->

### HTTP Statistics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| HTTP Server Requests Rate | Shows per-second rate of HTTP server requests count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| HTTP Server Errors Rate | Shows the count of HTTP server errors | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Duration Server Request | Shows per-second rate of the average time per server request | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| HTTP Client Requests Rate | Shows per-second rate of  HTTP Client requests count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| HTTP Client Errors Rate | Shows the count of HTTP Client errors | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Duration Client Request | Shows per-second rate of the average time per client request | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Spring: HTTP Connections Pool | Shows Tomcat and Jetty<br/> - busy thereads<br/> - current threads<br/> - max threads | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Quarkus: HTTP Server Worker Pool | Shows connections pool for Quarkus applications | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### gRPC Statistics (Quarkus)

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| gRPC Server Requests Received / Responses Sent | Shows per-second rate of gRPC Server received requests and responses sent count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| gRPC Server Unary Call Error Rate | Shows the count of gRPC Server errors |  |  |
| gRPC Server Average Unary Call Processing Duration | Shows per-second rate of the average time of gRPC call processing by server | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| gRPC Client Responses Received / Requests Sent | Shows per-second rate of gRPC Server received responses and sent requests count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| gRPC Client Unary Call Error Rate | Shows the count of gRPC Client errors |  |  |
| gRPC Client Average Unary Call Processing Duration | Shows per-second rate of the average time of gRPC call processing by client | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### gRPC Statistics (Spring)

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| gRPC Server Unary Requests by OK Code | Shows amount of successful requests to different methods per second for gRPC Server | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| gRPC Server Unary Failed Requests | Shows amount of error requests to different methods per second partitioned by error codes |  |  |
| gRPC Server 95%-tile Latency of Unary Response Time | How long it took to process an unary request, partitioned by method | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| gRPC Client Unary Requests by OK Code | Shows amount of successful requests to different methods per second for gRPC Client | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| gRPC Client Unary Failed Requests | Shows amount of error requests to different methods per second partitioned by error codes for gRPC Client |  |  |
| gRPC Client 95%-tile Latency of Unary Response Time | How long it took to process an unary request, partitioned by method for gRPC Client | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Garbage Collection

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Collections | Shows per-second rate of GC pause increase | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Pause Durations | Shows per-second rate of average increasing pause time per  pause count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Allocated/Promoted | Shows per-second rate of increasing memory allocated and promouted by GC | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### JDBC/Hikari

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| JDBC Connections | Show current count of JDBC connections and them limits | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Hickari Pool: Connections | Show count of connections via Hickari poll by statuses  | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Hickari Pool: Usage | Show statistic of using Hickari Pool connections, by time and by connections count | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Agroal

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Connection Pool | Shows connection pool for Agroal | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Created and Destroyed Connections | Shows a number of created and destroyed connections | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Connections Timings | Shows average time to create connections and acquire them for Agroal |  |  |
| Acquired Rate | Shows acquired operations per second | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Leak Detected | Shows number of times a leak was detected (a single connection can be detected multiple times) |  |  |
<!-- markdownlint-enable line-length -->

### Cassandra Sessions

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Average CQL Response Time | How long Cassandra process a CQL request (average time) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| 95%-tile Latency of CQL Response Time | How long Cassandra process a CQL request (95th percentile) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| 99%-tile Latency of CQL Response Time | How long Cassandra process a CQL request (99th percentile) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Receive/Sent Data | Shows network traffic in bytes per second for Cassandra sessions | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| CQL Timeouts Rate | Shows count of CQL timeouts | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| CQL Request Errors Rate | Shows count of errors during CQL requests | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Connection Errors Rate | Shows count of errors during connection to Cassandra | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Throttling Errors Rate | Shows count of throttling errors | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### RabbitMQ Client

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Connections | Shows number of current connections from the service to RabbitMQ | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Channels | Shows number of current RabbitMQ channels | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Published Messages Rate | Shows per-second rate of published messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Acknowledged Messages Rate | Shows per-second rate of acknowledged messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Unrouted Published Messages Rate | Shows per-second rate of unrouted published messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Unacknowledged Published Messages Rate | Shows per-second rate of unacknowledged published messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Acknowledged Published Messages Rate | Shows per-second rate of acknowledged published messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Consumed Messages Rate | Shows per-second rate of consumed messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Failed to Publish Messages Rate | Shows per-second rate of failed to publish messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Rejected Messages Rate | Shows per-second rate of rejected messages | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Classloading

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Classes loaded | Shows count of classes loaded by class loader | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Class delta | Shows delta of loaded classes count  | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Buffer Pools

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Direct Buffers Size | Shows buffer size<br/> - used<br/> - total available | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Direct Buffers Count | Shows estimation of the number of buffers in the pool | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Mapped Buffers Size | Shows estimation of <br/> - the memory that the Java virtual machine is using for this buffer pool<br/> - capacity of the buffers in this pool | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Mapped Buffers Count | Shows estimation of the number of buffers in the pool | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
