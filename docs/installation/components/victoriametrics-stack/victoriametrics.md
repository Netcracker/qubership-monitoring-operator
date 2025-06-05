#### VictoriaMetrics

For integration with VictoriaMetrics into Monitoring, the chart has the ability to specify RemoteWrite settings.

The following image illustrates how the metrics are sent during the use of this service in an external server.

![Prometheus to VictoriaMetrics](../../../images/prometheus_k8s_integration_victoriametrics.png)

To setup this integration, you need to configure:

* The Remote write section into Prometheus. It needs to send the metrics from Prometheus to VictoriaMetrics.

For example:

```yaml
prometheus:
  externalLabels:
    cluster: my_awesome_cluster
  remoteWrite:
    - url: "http://<victoriametrics_addr>:8428/api/v1/write"
      queueConfig:
        maxSamplesPerSend: 10000
        capacity: 20000
        maxShards: 30
```

The full list of parameters can be viewed at:

* [Prometheus](#prometheus)


