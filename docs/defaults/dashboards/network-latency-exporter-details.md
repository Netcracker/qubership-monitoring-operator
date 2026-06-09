# Network Latency Details

The dashboard shows network RTT metrics per source host / destination host pairs

## Tags

* `network`
* `k8s`

## Panels

### From $source to $dest

**Row From $source to $dest is multiplied by parameter `dest`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Average Mean RTT With Packet Loss | Shows average mean Round Trip Time (RTT) metrics over selected time range (with package loss). Displayed values is a values for whole probe which includes $packets_num packets sent |  |  |
| Packet Loss | Shows percent of packet loss per probe (which includes $packets_num packets sent) |  |  |
| Round Trip Time (RTT) | Show Round Trip Time (RTT) average time per probe (which includes $packets_num packets sent) |  |  |
| RTT Deviation | Shows standard deviation of Round Trip Time (RTT)  mean value per probe (which includes $packets_num packets sent) |  |  |
<!-- markdownlint-enable line-length -->
