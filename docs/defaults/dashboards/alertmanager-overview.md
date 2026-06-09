# Alertmanager Overview

Dashboard showing Prometheus Alertmanager metrics for observing status of the cluster and possible debugging.

## Tags

* `alertmanager`
* `self-monitor`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of instances | Number of Alertmanager instances | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Cluster size | Number of peers in the Alertmanager cluster. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Instance versions and up time | Table containing list of Alertmanager instances showing it's version, up time, last reload time and if it was successful. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Number of active alerts | Current number of active alerts. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Number of suppressed alerts | Current number of suppressed alerts. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Number of active silences | Current number of active silences. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Notifications

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Notifications sent from $instance | Number of sent notifications to distinct integrations such as PagerDuty, Slack and so on. On negative axis are displayed failed notifications. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Notification durations per integration on $instance | Duration of notification sends in 0.99 and 0.9 quantiles per integration. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
<!-- markdownlint-enable line-length -->

### Alerts

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Active alerts in $instance | Number of alerts by state such as `active`, `suppressed` etc. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Received alerts by status for $instance | Number of received alerts from Prometheus by status `firing` on positive axis and `resolved` on negative axis. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
<!-- markdownlint-enable line-length -->

### Cluster members

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Cluster health score for $instance | Shows cluster score representing cluster health. From Hashicorps official documentation: <br/>> This metric describes a node's perception of its own health based on how well it is meeting the soft real-time requirements of the protocol. This metric ranges from 0 to 8, where 0 indicates "totally healthy".<br/><br/>For more info see <https://www.consul.io/docs/agent/telemetry>.html#cluster-health | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Cluster members count on $instance | Shows gossip cluster members count in time and failing peers in case of any in red color. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Cluster peers left/joined on $instance | On positive axis shows number of peers that joined the cluster and on negative axis number of peers that left the cluster. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Cluster reconnections on $instance | On positive axis is number of attempts to reconnect the cluster. On negative axis if number of failed attempts. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Cluster messages count on $instance | On positive axis is number of sent cluster messages by type `update`  or `full_state` and on negative axis the same for received messages. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Cluster messages size on $instance | On positive axis is size of sent cluster messages by type `update`  or `full_state` and on negative axis the same for received messages. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Cluster messages queue on $instance | On positive axis is number of queued cluster messages and on negative axis number of pruned messages. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
<!-- markdownlint-enable line-length -->

### Gossip messages

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Count of oversized gossip messages on $instance | Number of oversized gossip message sent by $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Number of  propagated gossip messages on $instance | Number of  propagated gossip messages on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Duration of oversized gossip messages on $instance | Duration of oversized gossip message requests on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
<!-- markdownlint-enable line-length -->

### Nflog

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Nf log queries count for $instance | Number of log queries for $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Nf log query duration for $instance | Nf log query duration for $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Nf log snapshot size for $instance | Snapshot size for NF log on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Nf log Go GC time for $instance | Duration of the last notification log garbage collection cycle for $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Nf log snapshot duration for $instance | The duration of creating snapshoot fo NF log on  $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
<!-- markdownlint-enable line-length -->

### Silences

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Silences count by state on $instance | Number of silences by state on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Silences query count on $instance | Number of silences queries and errors on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Silences query duration on $instance | Silences query duration on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Silences snapshot size on $instance | Size of the silence snapshot in bytes on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Silences snapshot duration on $instance | Duration of the silence snapshot on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
| Silences GC duraton on $instance | Duration of the silence garbage collection cycle<br/> on $instance | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `instance`** |
<!-- markdownlint-enable line-length -->
