# Network Latency Overview

The dashboard shows overall network RTT metrics for all sources and destinations

## Tags

* `network`
* `k8s`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Average Ping Measurements | Shows average ping measurements grouped by destination host over selected time range | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Source: $source

**Row Source: $source is multiplied by parameter `source`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Average Mean RTT | Shows average mean RTT metrics over selected time range (with package loss). Displayed values is a values for whole probe which includes $packets_num packets sent |  |  |
| Ping Measurements From $source To All Hosts | Shows average values of ping measurements grouped by destination host over selected time range. Click on a destination node to open dashboard with detailed information | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
