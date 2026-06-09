# Kubernetes / Top Resources

Show first N (can be selected) resources by CPU, Memory, Disk usage

## Tags

* `k8s`
* `top`

## Panels

### Top pods by CPU usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by CPU usage | Show first $topN pods by CPU consumption.<br/>Show CPU usage in millicores, where `1000 millicores = 1 core`. In legend pods show as `namespace/pod`. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN by CPU usage | Show first $topN namespaces by CPU consumption.<br/>Show CPU usage as:<br/>* CPU usage in millicores per pods<br/>* CPU requests per pods<br/>* CPU limits per pods | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top pods by Memory usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by Memory usage | Show first $topN pods by Memory consumption.<br/>In legend pods show as `namespace/pod` | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN by Memory usage | Show first $topN pods by Memory consumption.<br/>Show Memory usage as:<br/>* Memory usage per pod<br/>* Memory requests per pod<br/>* Memory limits per pod | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top pods by Disk usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by Disk usage (Reads and Writes) | Show first $topN pods by Disk consumption. In legend pods show as `namespace/pod` | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN by Disk usage (Reads and Writes) | Show first $topN pods by Disk writes/reads consumption.<br/>Show Disk usage as:<br/><br/>* Disk reads per seconds<br/>* Disk writes per seconds | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top throttled pods

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN throttled pods, rate by seconds | Show first $topN pods which were throttled by CPU limits.<br/><br/>Show as rate by `container_cpu_cfs_throttled_seconds_total`. It means that this graph should show how many time pod spend in throttling per second. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top namespaces by CPU usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by CPU usage | Show first $topN namespace by CPU consumption.<br/>Show CPU usage in millicores, where `1000 millicores = 1 core`. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN by CPU usage | Show first $topN namespaces by CPU consumption.<br/>Show CPU usage as:<br/>* Sum of CPU usage in millicores in specified namespace<br/>* Sum of CPU requests for all pods in specified namespace<br/>* Sum of CPU limits for all pods in specified namespace | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top namespace by Memory usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by Memory usage | Show first $topN namespace by Memory consumption. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN by Memory usage | Show first $topN namespaces by Memory consumption.<br/>Show Memory usage as:<br/>* Sum of Memory usage in millicores in specified namespace<br/>* Sum of Memory requests for all pods in specified namespace<br/>* Sum of Memory limits for all pods in specified namespace | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top namespaces by Disk usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by Disk usage (Reads and Writes) | Show first $topN namespaces by Disk consumption | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN by Disk usage (Reads and Writes) | Show first $topN namespaces by Disk writes/reads consumption.<br/>Show Disk usage as:<br/><br/>* Disk reads per seconds<br/>* Disk writes per seconds | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top nodes by CPU usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by CPU usage, % | Show first $topN Kubernetes/OpenShift nodes by CPU consumption.<br/>Show CPU usage as CPU total usage percent (CPU usage / CPU cores). <br/><br/>Because there are a lot CPU spaces this metric calculate as:<br/><br/>`(1 - idle usage) * 100`<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN nodes by CPU usage | Show first $topN Kubernetes/OpenShift nodes by CPU consumption.<br/>Show CPU usage as:<br/>* CPU cores per nodes<br/>* CPU Load Average to last 5 minutes<br/>* CPU total usage percent (CPU usage / CPU cores)<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top nodes by Memory usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN nodes by Memory usage, % | Show first $topN Kubernetes/OpenShift nodes by Memory consumption.<br/>Show Memory usage in percents from total node memory.<br/><br/>As base metrics use `node_memory_MemAvailable_bytes` to show used bytes per node<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN nodes by Memory usage | Show first $topN Kubernetes/OpenShift nodes by Memory consumption.<br/>Show Memory usage as:<br/><br/>* Memory on node<br/>* Memory usage percent<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Top nodes by Disk usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Top $topN by Disk usage (Reads and Writes) | Show first $topN Kubernetes/OpenShift nodes by Disk consumption<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Top $topN nodes by Disk usage | Show first $topN Kubernetes/OpenShift nodes by Disk writes/reads consumption.<br/>Show Disk usage as:<br/><br/>* Disk reads per seconds<br/>* Disk writes per seconds<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
