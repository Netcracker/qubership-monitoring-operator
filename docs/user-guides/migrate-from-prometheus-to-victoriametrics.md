This guide describes how to migrate from Prometheus to VictoriaMetrics.

# Overview

Monitoring can work with Prometheus or VictoriaMetrics stack. It is two different TSDB and you can choose between them.
We recommend use VictoriaMetrics instead of Prometheus since monitoring-operator `v0.55`. If you used Prometheus before
you can deploy monitoring with VictoriaMetrics without problems.

Migration from Prometheus to VictoriaMetrics can be realized with two ways:

1. Without keeping the metrics
2. With saving the metrics

**NOTE:** Make sure that you need to keep the metrics from Prometheus. If your Prometheus has retention period 24h
(or 2-3 days) in emptyDir you do not need to migrate data to VictoriaMetrics, but if retention 7d(at least) it can be
the right way for saving data.

## Without keeping metrics

During deploy you should disable Prometheus installation and enable VictoriaMetrics stack installation. After successful
installation components of VictoriaMetrics will be created instead of Prometheus instance.

Prometheus stack is:

* `Prometheus-operator`
* `Prometheus`
* `AlertManager`

Instead of it you should use VictoriaMetrics stack. It contains:

* `VmOperator`
* `VmSingle` - TSDB
* `VmAgent` - Agent for scrape metrics
* `VmAlertManager` - VictoriaMetrics alerting system

**NOTE:** VictoriaMetrics has more components(vmOperator, vmSingle, vmAgent, vmAlertManager, vmAlert, vmAuth and
vmUser) than Prometheus. You do not have to enable installation for all of it. Choose just what you need.

Example:

```yaml
....
prometheus:
  install: false
    
alertManager:
  install: false

....
victorimetrics:
  vmOperator:
    install: true
  vmSingle:
    install: true
  vmAgent:
    install: true
  vmAlertManager:
    install: false
  vmAlert:
    install: false
....
```

**NOTE:** Description for all parameters [see here](../installation/README.md#victoria-metrics)

## Save metrics

### Before you begin

* Retention period in Prometheus at least 7 days (or even more)
* VictoriaMetrics tool - `vmctl`. For more information [see here.](https://docs.victoriametrics.com/vmctl.html#migrating-data-from-prometheus)
**NOTE:** `vmctl` doesn't support self-signed and signed internal CAs certificates.
* If you use ingress to vmsingle for uploading data you will be able to face with error `413 Request Entity Too Large`.
To resolve it ingress should have annotation `nginx.ingress.kubernetes.io/proxy-body-size: 100m`. Blocks of Prometheus
data can be bigger than ingress allows by default.

### Steps for data migration

For saving important metrics from Prometheus you should make manual steps.

Step 1. Make Prometheus snapshot

TSDB Admin APIs allows doing snapshots. You should to send POST request on Prometheus endpoint
`/api/v1/admin/tsdb/snapshot`. The request can be sent from local host via Kubernetes ingress
(connected with Prometheus) or from Kubernetes pod that has a tool that allows to do POST request (e.g. curl).

```shell
curl -XPOST http://<your_prometheus_URL>/api/v1/admin/tsdb/snapshot
{
  "status": "success",
  "data": {
    "name": "20171210T211224Z-2be650b6d019eb54"
  }
}
```

Endpoint has optional parameter `skip_head=<bool>` - skip data present in the head block.

**NOTE**: The snapshot is stored at `/prometheus/snapshots/` into Prometheus pod.

Step 2. Download on local host/another storage in Kubernetes

When snapshot is done you should to save on local host or another storage in Kubernetes. It needs to be done because
during update from Prometheus to VictoriaMetrics all Prometheus components will be deleted and data will be cleaned.

For copying snapshots on local host you can use `cp` command in `kubectl` tool:

```shell
kubectl cp -n <platform-monitoring-namespace> <prometheus-pod-name>:<prometheus-snapshot-path> <path-on-your-host>
```

For example:

```shell
kubectl cp -n monitoring prometheus-k8s-0:/prometheus/snapshots/20231025T141722Z-734900232b68c299 ./snapshot_1
```

Step 3. Deploy monitoring

The same as installation [without keeping metrics](#without-keeping-metrics). You should disable Prometheus installation
and enable VictoriaMetrics components. But you need to pay attention on `retentionPeriod` parameter. If it is less than
the first timestamp of metrics from Prometheus you will get error about it.

```yaml
....
prometheus:
  install: false

alertManager:
  install: false
  
....
victorimetrics:
  vmOperator:
    install: true
  vmSingle:
    install: true
    # Retention period specifies how long metrics will be stored. Timestamp of metrics in Prometheus snapshot have to
    # present between `now` and `now - retentionPeriod`.
    retentionPeriod: "14d"
    # Via this ingress snapshots will be uploaded in vmSingle from local host.
    ingress:
      install: true
      host: vmsingle.<cloud>.com
  vmAgent:
    install: true
  vmAlertManager:
    install: false
  vmAlert:
    install: false
....
```

**NOTE:** VictoriaMetrics has more components(vmOperator, vmSingle, vmAgent, vmAlertManager, vmAlert, vmAuth and
vmUser) than Prometheus. You do not have to enable installation for all of it. Choose just what you need. Description
for all parameters [see here](../installation/README.md#victoria-metrics)

Step 4. Run `vmctl` tool

After successful deploy you can import snapshots to the TSDB. VictoriaMetrics has `vmctl` tool for it. More information
[here](https://docs.victoriametrics.com/vmctl.html#migrating-data-from-prometheus).

* You can run tool locally:

  ```shell
  ./vmctl-prod prometheus --vm-addr=<addr-to-vmsingle> --prom-snapshot=<path-to-snapshots>
  ```

* Or you can create deployment in Kubernetes and use pod for import backup to VictoriaMetrics:

  ```yaml
  kind: Deployment
  apiVersion: apps/v1
  metadata:
    name: vmctl
    labels:
      app: vmctl
  spec:
    selector:
      matchLabels:
        app: vmctl
    template:
      metadata:
        creationTimestamp: null
        labels:
          app: vmctl
      spec:
        containers:
          - name: vmctl
            image: <vmctl-image> # e.g. victoriametrics/vmctl:v1.93.5
            command:
              - /bin/sh
              - '-c'
              - '--'
            args:
              - while true; do sleep 30; done;
  ```
  
  **NOTE**: You should mount or upload Prometheus backup to pod with vmctl tool.

The `vmctl` tool has a lot of flags it is just main parameters. But you can filter data in snapshot and upload not all
metrics from Prometheus.

For example:

```shell
./vmctl prometheus --prom-snapshot=/path/to/snapshot \
  --vm-addr=vmsingle.test_cloud.com \

Prometheus import mode
Prometheus snapshot stats:
  blocks found: 14;
  blocks skipped: 0;
  min time: 1581288163058 (2023/08/22T22:42:43Z);
  max time: 1582409128139 (2023/08/23T22:05:28Z);
  samples: 32549106;
  series: 27289.
Found 14 blocks to import. Continue? [Y/n] y

14 / 14 [-------------------------------------------------------------------------------------------] 100.00% 0 p/s
2023/08/23 15:50:03 Import finished!
2023/08/23 15:50:03 VictoriaMetrics importer stats:
  idle duration: 6.152953029s;
  time spent while importing: 44.908522491s;
  total samples: 32549106;
  samples/s: 724786.84;
  total bytes: 669.1 MB;
  bytes/s: 14.9 MB;
  import requests: 323;
  import requests retries: 0;
2023/08/23 15:50:03 Total time: 51.077451066s
```

Step 5. Check data in VictoriaMetrics

   1. Open vmui (Web UI for VictoriaMetrics). Created ingress from previous step can be used for it.

   2. Choose time range based on period of backup(e.g. Last 7 days).

   3. Execute simple query. It should contain metric from Prometheus backup. For example, you can use
      `container_cpu_system_seconds_total` metric.

   4. You should see list of metrics.
