This document provides information about integration with IBM Netcool alerting.

# IBM Netcool Operations Insight

IBM Netcool Network Management helps communication service providers, enterprise data centers,
and networking staff to discover, visualize, detect, configure, activate, integrate,
and remediate your network.

Supported features matrix:

| **Monitoring System**     | **Metrics** | **Dashboards** | **Alerting** | **Autoscaling** |
| ------------------------- | ----------- | -------------- | ------------ | --------------- |
| IBM Netcool               | ?           | ?              | ✓ Yes        | ?               |

**Legend**:

* `Yes` - feature supported and implemented
* `No` - feature not supported by target monitoring system
* `?` - not investigate
* `-` - not applicable

The integration between IBM Netcool and Monitoring Stack can be configured:

* for Monitoring in Cloud -  **after** deploying Monitoring and doesn't require any updates
* for monitoring in the VM - **during** deploy on the VM

# Integration with IBM Netcool

This section describes how to configure integration with IBM Netcool for Alertmanager deployed
in different places.

## Before you begin

* Find the address of IBM Netcool that provides the webhook integration.
* Check that Monitoring has access to the address.
* Check that you have permissions to `get`, `list`, `create`, `update`, and maybe `delete`.
  Custom Resources from the `monitoring.coreos.com/v1alpha1` group with kind `AlertmanagerConfig`
  (for Alertmanager in the Cloud).
* Determine which alerts you want to send in IBM Netcool and check that all of them have
  any specific signs (like labels and part of the name) to specify them in the selector.

Example of necessary permissions:

```yaml
- apiGroups:
    - "monitoring.coreos.com"
  resources:
    - alertmanagerconfigs
  verbs:
    - 'create'
    - 'delete'
    - 'get'
    - 'list'
    - 'update'
    - 'watch'
```

## Integration Alertmanager in the Cloud

Prometheus can be integrated with IBM Netcool by alerts. It means that Alertmanager,
as a part that is deployed with Prometheus, can send alerts to IBM Netcool.

The official documents of IBM Netcool provide information about how to configure it when Prometheus
uses various ConfigMaps to store its configurations.

However, in the current solution, prometheus-operator is used to deploy and configure Prometheus and
Alertmanager. So this guide describes how to configure the integration using approaches applicable
to prometheus-operator.

### Configure Monitoring

To configure Monitoring for integration, create an AlertmanagerConfig Custom Resource (CR).

**Note**: The `send_resolved` flag should be set to `true` so that the probe receives
resolution events.

1. Take a template and save it to the file, for example, with the name `alertmanager-netcool-webhook.yaml`:

    ```yaml
    apiVersion: monitoring.coreos.com/v1alpha1
    kind: AlertmanagerConfig
    metadata:
      name: alertmanager-netcool-webhook
      labels:
        app.kubernetes.io/component: monitoring  # Mandatory label
    spec:
      route:
        groupWait: 10s
        groupInterval: 5m
        repeatInterval: 3h
        receiver: 'netcool_probe'
    receivers:
    - name: 'netcool_probe'
        webhookConfig:
        - url: 'http://<ip_address>:<port>/probe/webhook/prometheus'
          sendResolved: true
    ```

2. Replace `<ip_address>:<port>` on the address and port for your environment.
3. Change the group condition as required.
4. Save the file and create it in Cloud:

   ```bash
   kubectl apply -f alertmanager-netcool-webhook.yaml
   ```

Prometheus usually takes a couple of minutes to reload the updated configuration and apply
the new configuration.
5. Verify that Prometheus events appear on the `OMNIbus Event List`.

## Integration Alertmanager on the VM

To configure Alertmanager on the VM that deployed as a part of External System Monitoring (ESM) for integration,
use `alertmanager_receivers` and `alertmanager_route` parameters during deployment.

**Note**: The `send_resolved` flag should be set to `true` so that the probe receives
resolution events.

1. Update ESM inventory file used for installation:

    ```yaml
    global:
      ...
    nodes:
    - name: ...
        ...
      role: "monitoring_server"
    roles:
      monitoring_server:
        services:
        - alertmanager
    services:
      ...
      alertmanager:
        install: true
        ...
        receivers:
        - name: 'netcool_probe'
          webhookConfig:
          - url: 'http://<ip_address>:<port>/probe/webhook/prometheus'
            sendResolved: true
        - name: blackhole
        route:
          continue: false
          receiver: blackhole
          group_wait: 10s
          group_interval: 30s
          repeat_interval: 1h
          routes:
            - continue: true
              matchers:
                - name: "<enable_notification_flag>"
                  value: "true"
                  matchType: "="
              receiver: netcool_probe
    ```

2. Replace matchers to project specific rules. In this example, `<enable_notification_flag>` labels
   what is used by the project to enable notifications for alerts.
3. Replace `<ip_address>:<port>` on the address and port for your environment.
4. Change the group condition as required.
5. Update inventory file and run update job.
   Prometheus AlertManager usually takes a couple of minutes to reload the updated configuration
   and apply the new configuration.
6. Verify that Prometheus AlertManager events appear on the `OMNIbus Event List`.

# Links

* Official documentation about integration IBM Netcool and Prometheus -
  [https://www.ibm.com/docs/en/netcoolomnibus/8?topic=cc-integrating-prometheus-alert-manager-netcool-operations-insight](https://www.ibm.com/docs/en/netcoolomnibus/8?topic=cc-integrating-prometheus-alert-manager-netcool-operations-insight)
* IBM Netcool OMNIBus -
  [https://www.ibm.com/docs/en/netcoolomnibus/8.1?topic=overview-introduction-tivoli-netcoolomnibus](https://www.ibm.com/docs/en/netcoolomnibus/8.1?topic=overview-introduction-tivoli-netcoolomnibus)
