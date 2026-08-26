# Yace Exporter Configuration Examples

This section contains ready-to-use YACE (Yet Another CloudWatch Exporter) configurations. Paste the content of a file
into the `yaceExporter.config` parameter of your deploy parameters.

## Available Configurations

| Configuration File | Description |
|--------------------|-------------|
| [config.yaml](config.yaml) | Discovery jobs for RDS, EBS, S3, Application ELB, and EFS in a single AWS account |
| [multi-account-multi-region.yaml](multi-account-multi-region.yaml) | Cross-account and multi-region collection through assumed IAM roles |

## Single Account

Discovery jobs find tagged resources and export the tags listed in `exportedTagsOnMetrics` as metric labels.
Discovery requires the `tag:GetResources` IAM permission.

```yaml title="config.yaml"
--8<-- "examples/components/yace-exporter-config/config.yaml"
```

## Multiple Accounts and Regions

One deployment covers every account and region. YACE assumes each role in `roles` and scrapes it in each region in
`regions`, so a job with two roles and two regions produces four scrape targets. The home identity must be allowed to
call `sts:AssumeRole` on each target role, and the trust policy of each target role must allow the home identity.

```yaml title="multi-account-multi-region.yaml"
--8<-- "examples/components/yace-exporter-config/multi-account-multi-region.yaml"
```

## Related Documentation

* [Yace-exporter component reference](../../../installation/components/exporters/yace-exporter.md)
* [AWS integration guide](../../../integration/amazon-aws.md#cloudwatch-metrics-with-yace-yet-another-cloudwatch-exporter)
* [Deploy parameter examples](../../deploy-parameters/yace-exporter)
