### graphite-remote-adapter

<!-- markdownlint-disable line-length -->
| Field                    | Description                                                                                                                                                                                                                | Scheme                                                                                                                       |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                  | Allows to enable or disable deploy graphite-remote-adapter.                                                                                                                                                                | *bool                                                                                                                        |
| name                     | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                     | string                                                                                                                       |
| image                    | A Docker image to deploy the graphite-remote-adapter.                                                                                                                                                                      | string                                                                                                                       |
| replicas                 | Number of created pods.                                                                                                                                                                                                    | int                                                                                                                          |
| resources                | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext          | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                         | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations              | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                    | []v1.Toleration                                                                                                              |
| nodeSelector             | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                         | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations              | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"     | map[string]string                                                                                                            |
| labels                   | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                           | map[string]string                                                                                                            |
| servicemonitor           | ServiceMonitor holds configuration attributes for graphite-remote-adapter.                                                                                                                                                 | object                                                                                                                       |
| servicemonitor.install   | Allows to disable create ServiceMonitor CR for graphite-remote-adapter during the deployment.                                                                                                                              | bool                                                                                                                         |
| servicePort              | The port for graphite-remote-adapter service.                                                                                                                                                                              | int                                                                                                                          |
| writeCarbonAddress       | The `host:port` of the Graphite server to send samples to.                                                                                                                                                                 | string                                                                                                                       |
| readUrl                  | The URL of the remote Graphite Web server to send samples to.                                                                                                                                                              | string                                                                                                                       |
| defaultPrefix            | The prefix to prepends to all metrics exported to Graphite.                                                                                                                                                                | string                                                                                                                       |
| enableTags               | Enable using Graphite tags.                                                                                                                                                                                                | string                                                                                                                       |
| priorityClassName        | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                      | string                                                                                                                       |
| additionalGraphiteConfig | Additional Graphite Config.                                                                                                                                                                                                | object                                                                                                                       |
| graphite                 | Graphite configuration                                                                                                                                                                                                     | object                                                                                                                       |
| write                    | Write graphite configuration configuration                                                                                                                                                                                 | object                                                                                                                       |
| compress_type            | Write graphite configuration configuration compress type. Supported values: "", plain, lz4                                                                                                                                 | string                                                                                                                       |
| lz4_preferences          | Parameters for lz4 streaming compression                                                                                                                                                                                   | object                                                                                                                       |
| compression_level        | LZ4 streaming compression level. Min value 3, max 12, default 9                                                                                                                                                            | int                                                                                                                          |
| auto_flush               | LZ4 streaming compression always flush; reduces usage of internal buffers. Default - false                                                                                                                                 | bool                                                                                                                         |
| decompression_speed      | LZ4 streaming compression parser favors decompression speed vs compression ratio. Works for high compression modes (compression_level >= 10) only. Default - false                                                         | bool                                                                                                                         |
| frame                    | Parameters for lz4 streaming compression frame.                                                                                                                                                                            | object                                                                                                                       |
| block_size               | The larger the block size, the (slightly) better the compression ratio. Larger blocks also increase memory usage on both compression and decompression sides. Values: max64KB, max256KB, max1MB, max4MB. Default: max64KB. | string                                                                                                                       |
| block_mode               | Linked blocks sharply reduce inefficiencies when using small blocks, they compress better. Default - false, i.e. disabled.                                                                                                 | bool                                                                                                                         |
| content_checksum         | Add a 32-bit checksum of frame's decompressed data. Default - false, i.e. disabled.                                                                                                                                        | bool                                                                                                                         |
| block_checksum           | Each block followed by a checksum of block's compressed data. Default - false, i.e. disabled.                                                                                                                              | bool                                                                                                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
graphite_remote_adapter:
  install: true
  name: graphite-remote-adapter
  image: product/prod.platform.system.graphite-remote-adapter:latest
  replicas: 1
  
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  servicemonitor:
    install: no
  servicePort: 9201
  writeCarbonAddress: localhost:9999
  readUrl: "http://guest:guest@localhost:8080"
  defaultPrefix: ""
  enableTags: true
  additionalGraphiteConfig:
    web:
      telemetry_path: "/metrics"
    write:
      timeout: 5m
    read:
      timeout: 5m
      delay: 1h
      ignore_error: true
    graphite:
      write:
        compress_type: lz4
        lz4_preferences:
          frame:
            block_size: max64KB
            block_mode: false
            content_checksum: false
            block_checksum: false
          compression_level: 9
          auto_flush: false
          decompression_speed: false
        carbon_transport: tcp
        carbon_reconnect_interval: 5m
        enable_paths_cache: true
        paths_cache_ttl: 4h
        paths_cache_purge_interval: 4h
        template_data:
          var1:
            foo: bar
          var2: foobar
```


