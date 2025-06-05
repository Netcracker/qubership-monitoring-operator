# CloudWatch Exporter Configuration Examples

This section contains ready-to-use CloudWatch Exporter configurations for various AWS services.

## Available Configurations

| Service | Configuration File | Description |
|---------|-------------------|-------------|
| **EBS** | [EBS.yaml](EBS.yaml) | Elastic Block Store volume metrics |
| **S3** | [S3.yaml](S3.yaml) | Simple Storage Service bucket and request metrics |
| **ELB** | [ELB.yaml](ELB.yaml) | Classic Load Balancer metrics |
| **Application ELB** | [ApplicationELB.yaml](ApplicationELB.yaml) | Application Load Balancer metrics |
| **Network ELB** | [NetworkELB.yaml](NetworkELB.yaml) | Network Load Balancer metrics |
| **EFS** | [EFS.yaml](EFS.yaml) | Elastic File System metrics |
| **Route53** | [Route53.yaml](Route53.yaml) | DNS service metrics |
| **Route53 Resolver** | [Route53Resolver.yaml](Route53Resolver.yaml) | DNS resolver metrics |
| **RabbitMQ** | [RabbitMQ.yaml](RabbitMQ.yaml) | Amazon MQ RabbitMQ metrics |
| **Cassandra** | [Cassandra.yaml](Cassandra.yaml) | Amazon Keyspaces (Cassandra) metrics |

## EBS Configuration

Monitor Amazon Elastic Block Store volumes for performance and utilization metrics.

```yaml title="EBS.yaml"
--8<-- "examples/components/cloudwatch-exporter-config/EBS.yaml"
```

### Key Metrics:
- **Volume I/O**: Read/Write bytes and operations
- **Performance**: Latency, throughput, queue length
- **Burst Credits**: For burstable volumes (gp2, st1, sc1)

## S3 Configuration

Monitor Amazon S3 buckets for storage usage and request metrics.

```yaml title="S3.yaml"
--8<-- "examples/components/cloudwatch-exporter-config/S3.yaml"
```

### Key Metrics:
- **Storage**: Bucket size and object count
- **Requests**: GET, PUT, DELETE, HEAD operations
- **Performance**: Request latency and error rates
- **Data Transfer**: Bytes uploaded/downloaded

## Usage

1. **Choose a configuration** from the table above
2. **Download the YAML file** by clicking on the link
3. **Apply the configuration** to your CloudWatch Exporter:

=== "Kubernetes"
    ```bash
    # Download and apply locally
    wget EBS.yaml
    kubectl apply -f EBS.yaml
    ```

=== "Helm Values"
    ```yaml
    # values.yaml
    cloudwatchExporter:
      config: |
        # Paste the content of chosen YAML file here
        region: us-east-1
        period_seconds: 120
        # ... rest of configuration
    ```

## Configuration Parameters

All configurations use these common parameters:

| Parameter | Default | Description |
|-----------|---------|-------------|
| `region` | `us-east-1` | AWS region to monitor |
| `period_seconds` | `120` | Metric collection interval |
| `delay_seconds` | `60` | Delay before collecting metrics |
| `aws_dimensions` | - | AWS resource dimensions for filtering |
| `aws_metric_name` | - | CloudWatch metric name |
| `aws_namespace` | - | AWS service namespace |
| `aws_statistics` | - | Aggregation functions (Average, Sum, etc.) |

## Customization

You can customize these configurations:

1. **Change regions**: Update the `region` field
2. **Adjust timing**: Modify `period_seconds` and `delay_seconds`
3. **Filter resources**: Add dimension filters
4. **Select metrics**: Remove unneeded metrics to reduce costs

## Related Documentation

- [CloudWatch Exporter Installation](../../../installation/components/exporters/cloudwatch-exporter.md)
- [AWS Integration Guide](../../../integration/amazon-aws.md)
- [Grafana Dashboards for AWS](grafana-dashboards-for-amazon-aws/)

## Grafana Dashboards

Pre-built Grafana dashboards are available for AWS services:

| Service | Dashboard File |
|---------|---------------|
| EBS | [ebs.json](grafana-dashboards-for-amazon-aws/ebs.json) |
| S3 | [s3.json](grafana-dashboards-for-amazon-aws/s3.json) |
| Application ELB | [application-elb.json](grafana-dashboards-for-amazon-aws/application-elb.json) |
| Classic ELB | [classic-elb.json](grafana-dashboards-for-amazon-aws/classic-elb.json) |
| Network ELB | [network-elb.json](grafana-dashboards-for-amazon-aws/network-elb.json) |
| EFS | [efs.json](grafana-dashboards-for-amazon-aws/efs.json) |
| RabbitMQ | [rabbitmq.json](grafana-dashboards-for-amazon-aws/rabbitmq.json) |
| Cassandra | [cassandra.json](grafana-dashboards-for-amazon-aws/cassandra.json) |

### Import Dashboards

1. Download the JSON file for your service
2. In Grafana, go to **+ → Import**
3. Upload the JSON file or paste its content
4. Configure data source if needed 