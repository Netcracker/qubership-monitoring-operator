This document describes which abilities monitoring has to collect metrics from services/microservices with endpoints
closed by TLS/mTLS (HTTPS).

# Overview

Let's describe how Prometheus (and VMAgent that can replace Prometheus in some cases) supports metrics collection
from endpoints with TLS/mTLS.

It supports:

| Type                        | Support? |
| --------------------------- | -------- |
| TLS                         | ✓        |
| mTLS                        | ✓        |
| skip certificate validation | ✓        |

To tell Prometheus that it should use an encrypted channel you need to specify it in `prometheus-operator`
Custom Resources (CRs) like `ServiceMonitor`, `PodMonitor` or `Probe`.

We are using the `prometheus-operator` and its CRs, and they contain all settings which tell Prometheus how to
find microservice, and how to scrape metrics from it. Also, these CRs contains information about TLS settings, like
references to certificates, server name and so on. Because all sensitive information in Kubernetes should store
in Secrets (or in other trusted storages like Vault) `ServiceMonitor`, `PodMonitor`, `Probe` allow specifying
only references to Secrets (or ConfigMap for CA and Cert certificates).

If you specified all parameters correctly, so the `prometheus-operator` with Prometheus will automatically
track your certificates, watch changes in them and use new (if the certificate is updated).

## Look under the hood

This section is not mandatory to read but may be useful to investigate issues or just to understand how the
`prometheus-operator` discovers Secrets and how it provides Secrets content to Prometheus.

When you specify the reference on the Secret in your `ServiceMonitor`, `PodMonitor` or `Probe` these CRs will
discover by the `prometheus-operator`. It watches all CRs and gets events when somebody adds, updates or removes CRs.

When the `prometheus-operator` discovers CRs it executes the following operations:

1. Get and read CR, for example, `ServiceMonitor`
2. Get and read Kubernetes Secret if CR config contains reference on the Secret
   1. Secret will read by name and only in the namespace where created CR, for example, `ServiceMonitor`
3. The `prometheus-operator` creates a new file in the Prometheus pod
   1. Content of Secrets or ConfigMaps on which referred CRs will create by path:

        ```yaml
        /etc/prometheus/certs/<type>_<namespace>_<object_name>_<key>
        ```

        for example:

        ```yaml
        /etc/prometheus/certs/secret_monitoring_kube-etcd-client-certs_etcd-client-ca.crt
        /etc/prometheus/certs/configmap_monitoring_examples-service_ca.crt
        ```

4. If Secrets or ConfigMaps will change the `prometheus-operator` will receive the event and update the content
   of the file in the Prometheus pod
5. Next, the `prometheus-operator` generates a config that will use local paths to certificate files. For example:

    ```yaml
    tls_config:
    ca_file: /etc/prometheus/certs/secret_monitoring_kube-etcd-client-certs_etcd-client-ca.crt
    cert_file: /etc/prometheus/certs/secret_monitoring_kube-etcd-client-certs_etcd-client.crt
    key_file: /etc/prometheus/certs/secret_monitoring_kube-etcd-client-certs_etcd-client.key
    ```

Of course, such automation leads to additional required permissions for `prometheus-operator`. It needs permissions
to list, get and watch Secrets and ConfigMap in all namespaces where it should work (or in cluster scope).

# How to use

This section describes how to configure TLS and mTLS to collect metrics from services.

## Requirements for microservices and restrictions

There are some requirements for microservices that they must comply with:

* Certificates **must** store in Kubernetes Secret. You can't use certificates inside the pods.
  * At least **CA** certificate **must** store in the Secret
* Secret **must** store in the namespace where configured `ServiceMonitor`, `PodMonitor`, `Probe`. You can't configure
  a reference on Secret in another namespace.
* TLS config **must** be configured for each endpoint that uses TLS encryption
* Other trusted storage, like Vault, are not supported now

## Configure TLS for scrape metrics

Under configure TLS we mean that for the configuration you will use only a CA certificate. Other certificates usually
need for mTLS configuration.

To configure metrics collection from the TLS endpoint for `ServiceMonitor` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        ca:
          secret:
            name: secret-name-with-certificate
            key: key-in-secret-with-certificate
```

Also if your CA certificate is stored in Kubernetes ConfigMap you can use the reference on it:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        ca:
          configMap:
            name: configMap-name-with-certificate
            key: key-in-configMap-with-certificate
```

To configure metrics collection from the TLS endpoint for `PodMonitor` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        ca:
          secret:
            name: secret-name-with-certificate
            key: key-in-secret-with-certificate
```

Also if your CA certificate is stored in Kubernetes ConfigMap you can use the reference on it:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        ca:
          configMap:
            name: configMap-name-with-certificate
            key: key-in-configMap-with-certificate
```

To configure metrics collection from the TLS endpoint for `Probe` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
...
spec:
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'https://example.com'
        - 'https://google.com'
    tlsConfig:                  # <-- Need to add this section
      ca:
        secret:
          name: secret-name-with-certificate
          key: key-in-secret-with-certificate
```

Also if your CA certificate is stored in Kubernetes ConfigMap you can use the reference on it:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
...
spec:
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'https://example.com'
        - 'https://google.com'
    tlsConfig:                  # <-- Need to add this section
      ca:
        configMap:
          name: configMap-name-with-certificate
          key: key-in-configMap-with-certificate
```

## Configure mTLS for scrape metrics

To configure mTLS you need to have all 3 certificates: CA certificate, key and cert certificates (private and public
parts of client certificate). Also in some cases, you may need to use `serverName` which is used to verify the hostname
for the targets.

To configure metrics collection from the TLS endpoint for `ServiceMonitor` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        ca:
          secret:
            name: secret-name-with-certificate
            key: key-with-ca-in-secret-with-certificate
        cert:
          secret:
            name: secret-name-with-certificate
            key: key-with-cert-in-secret-with-certificate
        keySecret:
          name: secret-name-with-certificate
          key: key-with-key-in-secret-with-certificate
```

Also if your CA certificate is stored in Kubernetes ConfigMap you can use the reference on it:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        ca:
          configMap:
            name: configMap-name-with-certificate
            key: key-in-configMap-with-certificate
        cert:
          configMap:
            name: configMap-name-with-certificate
            key: key-with-cert-in-configMap-with-certificate
        # For "keySecret" can't be configure reference on ConfigMap, only Secret allowed
        keySecret:
          name: secret-name-with-certificate
          key: key-with-key-in-secret-with-certificate
```

To configure metrics collection from the TLS endpoint for `PodMonitor` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- need to specify "https" value
      tlsConfig:     # <-- need to add this section
        ca:
          secret:
            name: secret-name-with-certificate
            key: key-in-secret-with-certificate
        cert:
          secret:
            name: secret-name-with-certificate
            key: key-with-cert-in-secret-with-certificate
        keySecret:
          name: secret-name-with-certificate
          key: key-with-key-in-secret-with-certificate
```

Also if your CA certificate is stored in Kubernetes ConfigMap you can use the reference on it:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- need to specify "https" value
      tlsConfig:     # <-- need to add this section
        ca:
          configMap:
            name: configMap-name-with-certificate
            key: key-with-ca-in-configMap-with-certificate
        cert:
          configMap:
            name: configMap-name-with-certificate
            key: key-with-cert-in-configMap-with-certificate
        # For "keySecret" can't be configure reference on ConfigMap, only Secret allowed
        keySecret:
          name: secret-name-with-certificate
          key: key-with-key-in-secret-with-certificate
```

To configure metrics collection from the TLS endpoint for `Probe` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
...
spec:
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'https://example.com'
        - 'https://google.com'
    tlsConfig:                  # <-- need to add this section
      ca:
        secret:
          name: secret-name-with-certificate
          key: key-with-ca-in-secret-with-certificate
      cert:
        secret:
          name: secret-name-with-certificate
          key: key-with-cert-in-secret-with-certificate
      keySecret:
        name: secret-name-with-certificate
        key: key-with-key-in-secret-with-certificate
```

Also if your CA certificate is store in Kubernetes ConfigMap you can use the reference on it:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
...
spec:
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'https://example.com'
        - 'https://google.com'
    tlsConfig:                  # <-- need to add this section
      ca:
        configMap:
          name: configMap-name-with-certificate
          key: key-with-ca-in-configMap-with-certificate
      cert:
        configMap:
          name: configMap-name-with-certificate
          key: key-with-cert-in-configMap-with-certificate
      # For "keySecret" can't be configure reference on ConfigMap, only Secret allowed
      keySecret:
        name: secret-name-with-certificate
        key: key-with-key-in-secret-with-certificate
```

## Configure skip certificate validation for scrape metrics from the TLS endpoint

**Warning!** Using TLS without certificate validation **strongly** is not recommended to use on production.

For cases when you have TLS endpoints and have no certificates, but want to configure metrics collection, you can
allow Prometheus to skip certificate validation.

To configure metrics collection from the TLS endpoint for `ServiceMonitor` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- Need to specify "https" value
      tlsConfig:     # <-- Need to add this section
        insecureSkipVerify: true
```

To configure metrics collection from the TLS endpoint for `PodMonitor` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that TLS config should be specify for each endpoint if it use TLS
    - interval: 30s
      path: /metrics
      scheme: https  # <-- need to specify "https" value
      tlsConfig:     # <-- need to add this section
        insecureSkipVerify: true
```

To configure metrics collection from the TLS endpoint for `Probe` need add the following configuration:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
...
spec:
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'https://example.com'
        - 'https://google.com'
    tlsConfig:                  # <-- need to add this section
      insecureSkipVerify: true
```

## How to configure TLS for additional-scrape-config

In very rare cases you may need to configure metrics collection using Prometheus configuration syntax.
For example, in cases when need to collect metrics from an endpoint outside the cloud.

And the `prometheus-operator` allows us to do it. In the namespace with monitoring exist a Secret with the name
`additional-scrape-config` which can add native Prometheus config.

I will explain how to configure TLS using job configuration to collect metrics from Graylog VM.

To do it need:

1. Add certificates in Kubernetes Secret, for example with the configuration:

    ```yaml
    apiVersion: v1
    kind: Secret
    metadata:
      name: example-service
    stringData:
      ca.crt: <certificate_content>
      certificate.key: <certificate_content>
      certificate.crt: <certificate_content>
   ```

2. Manually create this Secret in the namespace with the Monitoring
3. Add a link to Secret in Prometheus deploy parameters:

    ```yaml
    prometheus:
      secrets:
        - <secret_name_1>
        - <secret_name_2>
        ...
        - <secret_name_N>
    ```

4. The content of Secret will mount as files in the directory:

    ```yaml
    /etc/prometheus/secrets/<secret_name>/<key_name>
    ```

    for example:

    ```yaml
    /etc/prometheus/secrets/example-service/ca.crt
    /etc/prometheus/secrets/example-service/certificate.crt
    /etc/prometheus/secrets/example-service/certificate.key
    ```

5. Add config for `additional-scrape-config`. You can add it as deploy parameters for Monitoring:

    ```yaml
    - job_name: graylog
    honor_timestamps: true
    scrape_interval: 30s
    scrape_timeout: 10s
    metrics_path: /api/plugins/org.graylog.plugins.metrics.prometheus/metrics
    scheme: http
    static_configs:
      - targets:
        - 10.101.17.214
    basic_auth:
      username: admin
      password: <secret>
    tls_config:
      ca_file: /etc/prometheus/secrets/example-service/ca.crt
      cert_file: /etc/prometheus/secrets/example-service/certificate.crt
      key_file: /etc/prometheus/secrets/example-service/certificate.key
    ```

To configure metrics collection without certificate validation you need to configure only the parameter
`insecure_skip_verify: true`:

**Warning!** Using TLS without certificate validation **strongly** is not recommended to use on production.

```yaml
- job_name: graylog
  honor_timestamps: true
  scrape_interval: 30s
  scrape_timeout: 10s
  metrics_path: /api/plugins/org.graylog.plugins.metrics.prometheus/metrics
  scheme: http
  static_configs:
  - targets:
    - 10.101.17.214
  basic_auth:
    username: admin
    password: <secret>
  tls_config:
    insecure_skip_verify: true
```

# Examples
