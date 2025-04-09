his document describes how to configure TLS for Monitoring components and how to integrate them with the Cert-manager.

# Table of content

- [Table of content](#table-of-content)
- [Overview](#overview)
  - [Architecture](#architecture)
    - [Steps to renew certificates](#steps-to-renew-certificates)
    - [Restrictions and bugs](#restrictions-and-bugs)
  - [Parameters for Prometheus](#parameters-for-prometheus)
    - [webTLSConfig](#webtlsconfig)
      - [clientAuthType](#clientauthtype)
    - [generateCerts](#generatecerts)
    - [createSecret](#createsecret)
  - [Ingress configuration](#ingress-configuration)
  - [About cert-manager](#about-cert-manager)
    - [Monitoring-operator integration with cert-manager](#monitoring-operator-integration-with-cert-manager)
      - [Cert-manager certificates for Prometheus](#cert-manager-certificates-for-prometheus)
      - [Cert-manager certificates for AlertManager](#cert-manager-certificates-for-alertmanager)
      - [Cert-manager certificates for Grafana](#cert-manager-certificates-for-grafana)
      - [Cert-manager certificates for Victoriametrics](#cert-manager-certificates-for-victoriametrics)
  - [Examples](#examples)
    - [Example to create secret](#example-to-create-secret)
      - [Cert-manager certificates for Victoriametrics](#cert-manager-certificates-for-victoriametrics-1)
    - [Example with manually pre-created secret](#example-with-manually-pre-created-secret)
    - [Example full parameters with creating secret](#example-full-parameters-with-creating-secret)
    - [Simple example to use cert-manager](#simple-example-to-use-cert-manager)
    - [Full example to use cert-manager](#full-example-to-use-cert-manager)
    - [Example for prometheus-adapter-operator to use cert-manager](#example-for-prometheus-adapter-operator-to-use-cert-manager)
    - [Example for prometheus-adapter-operator to use createSecret](#example-for-prometheus-adapter-operator-to-use-createsecret)
    - [Example for prometheus-adapter-operator to use existingSecret](#example-for-prometheus-adapter-operator-to-use-existingsecret)
    - [Example for prometheus-adapter-operator to use multiple secrets containing certificates](#example-for-prometheus-adapter-operator-to-use-multiple-secrets-containing-certificates)
    - [Example for victoriametrics to use cert-manager](#example-for-victoriametrics-to-use-cert-manager)
    - [Example for victoriametrics to use createSecret](#example-for-victoriametrics-to-use-createsecret)
    - [Example for victoriametrics to use existingSecret](#example-for-victoriametrics-to-use-existingsecret)

# Overview

Prometheus supports Transport Layer Security (TLS) encryption for connections to Prometheus instances. If you would
like to enforce TLS for those connections, you would need to create a specific web configuration file.


[Back to TOC](#table-of-content)


## Architecture

The diagram below shows a simplified structure of the architecture and secure connections between components.

![Prometheus TLS architecture](../images/prometheus_k8s_tls.png)

In general, we provide ability to enable TLS in Prometheus,but connection between prometheus and
grafana\alertmanager is unsecured (HTTP) with option to skipping TLS verify.

If enable TLS in Prometheus we have secured connection between Prometheus and Ingress. Ingress must have certificates
to getting access to Prometheus UI.


[Back to TOC](#table-of-content)


### Steps to renew certificates

TLS certificates have a limited duration, and you can update them over time by replacing the content of
corresponding Secrets. Cert-manager can do that automatically if you set `renewBefore` parameter for certificate.

Unfortunately, in some cases, the renewal of certificates at destinations doesn't happen automatically and additional
steps are required. The table below reflects actions required after updating the certificates:

| Service            | No extra steps | Restart pod |
| ------------------ | -------------- | ----------- |
| Ingresses          | ✓              |             |
| Prometheus server  |                | ✓           |
| Prometheus targets | ✓              |             |
| AlertManager       | TBD            | TBD         |
| Grafana            | TBD            | TBD         |

After restart, Prometheus server may be unavailable for a short time.

**Warning!**: Notice that if Prometheus uses emptyDir as long-term storage, your **data will be lost** during restart.
To avoid this behavior, please, use PV or remote storage to store metrics.

You can set TLS configuration in Pod/Service Monitors with certificates. If you renew these certs,
prometheus-operator will apply changes to Prometheus targets configuration in runtime.


[Back to TOC](#table-of-content)


### Restrictions and bugs

* Now Grafana are getting metrics from Prometheus without certificate verify.
  Grafana allows enabling TLS, but certificates to store in open-view in GrafanaDataSource CRD. It is unsecure.
  [Feature request](https://github.com/grafana-operator/grafana-operator/issues/769)
  to Grafana community to add ability to store certificates in secret.
* Alertmanager doesn't support TLS.
  * [Feature request](https://github.com/prometheus-operator/prometheus-operator/issues/4830) to Prometheus
community to support TLS.
  * [Feature request](https://github.com/prometheus-operator/prometheus-operator/issues/4241) to Alertmanager supports mTLS.
  * [Pull request](https://github.com/prometheus-operator/prometheus-operator/pull/4593) in progress.


[Back to TOC](#table-of-content)


## Parameters for Prometheus

Section `prometheus.tlsConfig` defines the TLS parameters for Prometheus.

For more information, refer to
[https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#webtlsconfig](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#webtlsconfig)

<!-- markdownlint-disable line-length -->
| Field          | Description                                                                                                                                                                                                                                      | Scheme |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ |
| keySecret      | Secret containing the TLS key for the server. For more information, refer to *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretkeyselector-v1-core)                                             | object |
| cert           | Contains the TLS certificate for the server. For more information, refer to [SecretOrConfigMap](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#secretorconfigmap)                                     | object |
| client_ca      | Contains the CA certificate for client certificate authentication to the server. For more information, refer to [SecretOrConfigMap](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#secretorconfigmap) | object |
| clientAuthType | Server policy for client authentication. Maps to ClientAuth Policies. For more detail on clientAuth options: [https://golang.org/pkg/crypto/tls/#ClientAuthType](https://golang.org/pkg/crypto/tls/#ClientAuthType)                              | string |
| generateCerts  | Allows to configure generation of TLS certificate for Prometheus by [cert-manager](https://cert-manager.io/).                                                                                                                                    | object |
| createSecret   | Specifies content for secret that will be created.                                                                                                                                                                                               | object |
<!-- markdownlint-enable line-length -->


[Back to TOC](#table-of-content)


### webTLSConfig

Secret containing the TLS key for the server.

Fields `keySecret` and `cert` required for `webTLSConfig` section if you want to use pre-created Secrets.

Example:

```yaml
prometheus:
  tlsConfig:
    webTLSConfig:
      keySecret:
        name: prometheus-tls
        key: tls.key
      cert:
        secret:
          name: prometheus-tls
          key: tls.crt
      client_ca:
        secret:
          name: prometheus-tls
          key: tls.crt
```

#### clientAuthType

Server policy for client authentication. Maps to ClientAuth Policies.

Available values:

* `NoClientCert` indicates that no client certificate should be requested during the handshake, and if any certificates
  are sent they will not be verified.
* `RequestClientCert` indicates that a client certificate should be requested during the handshake, but does not
  require that the client send any certificates.
* `RequireAnyClientCert` indicates that a client certificate should be requested during the handshake, and that at
  least one certificate is required to be sent by the client, but that certificate is not required to be valid.
* `VerifyClientCertIfGiven` indicates that a client certificate should be requested during the handshake, but does
  not require that the client sends a certificate. If the client does send a certificate it is required to be valid.
* `RequireAndVerifyClientCert` indicates that a client certificate should be requested during the handshake, and that
  at least one valid certificate is required to be sent by the client.

**NOTE:** By default, `clientAuthType` is `VerifyClientCertIfGiven` if you will be use `client_ca`. `NoClientCert` is
unavailable with `client_ca`.

### generateCerts

Allows to configure generation of TLS certificate for Prometheus by [cert-manager](https://cert-manager.io/).

If `tlsConfig.generateCerts.enabled` is true, the monitoring-operator will create Issuer, Certificate and use it for
TLS in Prometheus.

If set to `true`, certs and key from generated secret instead of specified in the `tlsConfig.webTLSConfig` section.
Other parameters from the `tlsConfig.webTLSConfig` section (e.g. clientAuthType`) will not be overridden.

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                         | Scheme |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| enabled           | Enable integration with cert-manager.                                                                                                                                               | bool   |
| secretName        | Name of generated secret with certificate, CA certificate and private key. Default value is `prometheus-cert-manager-tls`.                                                          | string |
| clusterIssuerName | You can set this parameter to use Cluster Issuer. Otherwise, if this parameter is empty, self-signed non-cluster `Issuer` will be created and used in the new Certificate resource. | string |
| duration          | Defines duration of the certificate in days.                                                                                                                                        | int    |
| renewBefore       | Specifies how long before expiry a certificate should be renewed.                                                                                                                   | int    |
<!-- markdownlint-enable line-length -->

Example of all available parameters in this section:

```yaml
prometheus:
  tlsConfig:
    generateCerts:
      enabled: false
      secretName: prometheus-cert-manager-tls
      clusterIssuerName: ""
      duration: 365
      renewBefore: 15
```

### createSecret

Specifies content for secret that will be created.

<!-- markdownlint-disable line-length -->
| Field | Description                                                                                   | Scheme |
| ----- | --------------------------------------------------------------------------------------------- | ------ |
| name  | Name of secret with cert, ca and key                                                          | string |
| cert  | TLS certificate for prometheus. By default, cert will be stored in secret with `tls.crt` key. | string |
| key   | TLS key for prometheus. By default, cert will be stored in secret with `tls.key` key.         | string |
| ca    | TLS CA for prometheus. By default, cert will be stored in secret with `ca.crt` key.           | string |
<!-- markdownlint-enable line-length -->

Example of configuration:

```yaml
prometheus:
  tlsConfig:
    createSecret:
      name: web-prometheus-certificates
      cert: |-
        -----BEGIN CERTIFICATE-----
        ...
        -----END CERTIFICATE-----
      key: |-
        -----BEGIN PRIVATE KEY-----
        ...
        -----END PRIVATE KEY-----
      ca: |-
         -----BEGIN CERTIFICATE-----
        ...
        ------END CERTIFICATE-----
```

Example of secret that will be create:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: <prometheus.tlsConfig.secret.name>
  labels:
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/name: prometheus
data:
  ca.crt: <prometheus.tlsConfig.secret.ca>
  tls.crt: <prometheus.tlsConfig.secret.cert>
  tls.key: <prometheus.tlsConfig.secret.key>
```

## Ingress configuration

You can secure an Ingress by specifying a Secret that contains a TLS private key and certificate.
The Ingress resource only supports a single TLS port, 443.

**NOTE:** The TLS secret must contain keys named `tls.crt` and `tls.key` that contain the certificate and private
key to use for TLS. For example:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: testsecret-tls
  namespace: default
data:
  tls.crt: # base64 encoded cert
  tls.key: # base64 encoded key
type: kubernetes.io/tls
```

You can use `prometheus.ingress.tlsSecretName` for enabling TLS to Ingress. For example:

```yaml
...
prometheus:
  ...
  ingress:
    install: true
    host: prometheus.test.com
    tlsSecretName: prometheus-tls
  ...
...
```

**NOTE:** you can create and use secret that specified in `prometheus.tlsConfig.secret`.

For more information, refer to
[https://kubernetes.io/docs/concepts/services-networking/ingress/#tls](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls)


[Back to TOC](#table-of-content)


## About cert-manager

[Cert-manager](https://cert-manager.io/) is X.509 certificate controller for Kubernetes and OpenShift workloads. It will
obtain certificates from a variety of Issuers, both popular public Issuers as well as private Issuers, and ensure
the certificates are valid and up-to-date, and will attempt to renew certificates at a configured time before expiry.

In order for cert-manager to generate a secret containing certificates and private key, you need to take several steps:

0. Make sure cert-manager is installed on the cluster. It is usually installed in the `cert-manager` namespace.
1. You should use [Issuer or ClusterIssuer](https://cert-manager.io/docs/concepts/issuer/) resource for creating
   certificate. These resources represent certificate authorities (CAs) that are able to generate signed certificates
   by honoring certificate signing requests. Example of namespaced self-signed Issuer resource:

   ```yaml
   apiVersion: cert-manager.io/v1
   kind: Issuer
   metadata:
     name: prometheus-tls-issuer
   spec:
     selfSigned: {}
   ```

2. Then you can create [Certificate](https://cert-manager.io/docs/concepts/certificate/) resource. Configuration of this
   resource allows to change parameters of generated certificates and private key. You can find an example of
   certificate resource [here](https://cert-manager.io/docs/usage/certificate/#creating-certificate-resources).
3. Cert-manager will create [Certificate Request](https://cert-manager.io/docs/concepts/certificaterequest/)
   resource based on created Certificate resource.
4. Also, cert-manager will create Secret resource with name specified in the Certificate resource previously.
   Generated secret contains fields `ca.crt` with PEM CA certificate, `tls.crt` with PEM private key and `tls.key`
   with PEM signed certificate chain by default.

The generated secret can be used in pods later as a volume.


[Back to TOC](#table-of-content)


### Monitoring-operator integration with cert-manager

If you use monitoring-operator, you don't need to create cert-manager resources for its components manually.
These resources can be created by Helm or monitoring-operator during deploy in some cases.

#### Cert-manager certificates for Prometheus

You can specify [`prometheus.tlsConfig.generateCerts`](#generatecerts) section to use cert-manager for TLS in
Prometheus and the corresponding Ingress.

To create and use cert-manager resources, just specify the
`prometheus.tlsConfig.generateCerts.enabled` parameter. This is the only required
parameter in the section.

If you set `prometheus.tlsConfig.generateCerts.enabled` to `true` during deploy, Certificate CR and (optionally)
Issuer CR will be created in the namespaces with the monitoring-operator by Helm. Then cert-manager will make
Certificate Request CR and Secret with `ca.crt`, `tls.crt` and `tls.key` fields automatically. At the same time the
monitoring-operator will create Prometheus with the Secret generated by cert-manager that will be mounted as
volume. Generated certs will be used in enabled TLS automatically. Also, TLS with the same certs will be turned on in
the corresponding Prometheus' Ingress. _You do not need to specify parameters in the Ingress section in this case._

In addition to the `generateCerts.enabled` parameter the `generateCerts` section allows configuring
the following values:

* `generateCerts.secretName`: defines the name of the generated Secret (`prometheus-cert-manager-tls` by default).
* `generateCerts.clusterIssuerName`: allows to specify pre-created ClusterIssuer instead of Issuer created by Helm to
  verify generated certs. Notice that ClusterIssuer is a cluster-wide entity and should be created manually before
  deploy in any namespace on cluster. If `clusterIssuerName` is empty (by default), Helm will create self-signed Issuer.
* `generateCerts.duration`: allows configuring duration of generated certificates
  (integer value in **days**; `365` by default).
* `generateCerts.renewBefore`: specifies how long before expiry a certificate should be renewed
  (integer value in **days**; `15` by default).

**NOTE:** Prometheus cannot update certificates in the container automatically. See more in
[Steps to renew certificates](#steps-to-renew-certificates).

New certificates in Ingresses and Prometheus targets (these certs can be specified in Pod/Service monitors) reinstalled
automatically.


[Back to TOC](#table-of-content)


#### Cert-manager certificates for AlertManager

TBD (Monitoring-operator do not support TLS for AlertManager yet)


[Back to TOC](#table-of-content)


#### Cert-manager certificates for Grafana

TBD (Monitoring-operator do not support TLS for Grafana yet)


[Back to TOC](#table-of-content)


#### Cert-manager certificates for Victoriametrics

You can specify `victoriametrics.<vm_component>.tlsConfig.generateCerts` section to use cert-manager for TLS
in Victoriametrics. In order to enable TLS for all victoriametrics components, it is necessary to set
`victoriametrics.tlsEnabled` to `true`. By default, cert-manager will be used to generate certificates.
In order for cert-manager to generate certificates, it is necessary to provide the self-signed cluster issuer name
using `victoriametrics.clusterIssuerName`. The certificate CR and the secret with with `ca.crt`, `tls.crt` and
`tls.key` fields will be automatically generated by the operator using Helm and cert-manager respectively.

`generateCerts` section allows configuring the following values:

`generateCerts.secretName`: defines the name of the generated Secret (`<vm_component>-tls-secret` by default).

`generateCerts.duration`: allows configuring duration of generated certificates
(integer value in days; 365 by default).

`generateCerts.renewBefore`: specifies how long before expiry a certificate should be renewed
(integer value in days; 15 by default).


[Back to TOC](#table-of-content)


## Examples

### Example to create secret

```yaml
...
prometheus:
  tlsConfig:
    webTLSConfig:
      keySecret:
        name: prometheus-tls
        key: tls.key
      cert:
        secret:
          name: prometheus-tls
          key: tls.crt
    createSecret:
      name: prometheus-tls
      cert: |-
        -----BEGIN CERTIFICATE-----
        ...
        -----END CERTIFICATE-----
      key: |-
        -----BEGIN PRIVATE KEY-----
        ...
        -----END PRIVATE KEY-----
      ca: |-
        -----BEGIN CERTIFICATE-----
        ...
        -----END CERTIFICATE-----
...
```


[Back to TOC](#table-of-content)


#### Cert-manager certificates for Victoriametrics

You can specify `victoriametrics.<vm_component>.tlsConfig.generateCerts` section to use cert-manager for TLS
configuration in Victoriametrics.


[Back to TOC](#table-of-content)


### Example with manually pre-created secret

```yaml
...
prometheus:
  tlsConfig:
    webTLSConfig:
      keySecret:
        name: pre-created-secret
        key: tls.key
      cert:
        secret:
          name: pre-created-secret
          key: tls.crt
      client_ca:
        secret:
          name: pre-created-secret
          key: ca.crt
  ingress:
    install: true
    host: <host>
    tlsSecretName: pre-created-secret
...
```


[Back to TOC](#table-of-content)


### Example full parameters with creating secret

```yaml
...
prometheus:
  tlsConfig:
    webTLSConfig:
      clientAuthType: "VerifyClientCertIfGiven"
      keySecret:
        name: prometheus-tls
        key: tls.key
      cert:
        secret:
          name: prometheus-tls
          key: tls.crt
      client_ca:
        secret:
          name: prometheus-tls
          key: ca.crt
    createSecret:
      name: prometheus-tls
      cert: |
        <data>
      key: |
        <data>
      ca: |
        <data>
  ingress:
    install: true
    host: <host>
    tlsSecretName: prometheus-tls
...
```


[Back to TOC](#table-of-content)


### Simple example to use cert-manager

```yaml
...
prometheus:
  tlsConfig:
    generateCerts:
      enabled: true
...
```


[Back to TOC](#table-of-content)


### Full example to use cert-manager

```yaml
...
prometheus:
  tlsConfig:
    generateCerts:
      enabled: true
      secretName: prometheus-cert-manager-tls
      clusterIssuerName: ""
      duration: 365
      renewBefore: 15
...
```

### Example for prometheus-adapter-operator to use cert-manager

```yaml
...
prometheusAdapter:
  operator:
    ...
    tlsEnabled: true
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        clusterIssuerName: "dev-cluster-issuer"
        secretName: "prometheus-adapter-client-tls-secret"  
...
```

### Example for prometheus-adapter-operator to use createSecret

```yaml
...
prometheusAdapter:
  operator:
    ...
    tlsEnabled: true
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        secretName: "prometheus-adapter-client-tls-secret"
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        cert: |-
          -----BEGIN RSA PRIVATE KEY-----
          ...
          -----END RSA PRIVATE KEY-----
...
```

### Example for prometheus-adapter-operator to use existingSecret

```yaml
...
prometheusAdapter:
  operator:
    ...
    tlsEnabled: true
    tlsConfig:
      existingSecret: "prometheus-adapter-client-tls-secret"
...
```

### Example for prometheus-adapter-operator to use multiple secrets containing certificates

```yaml
...
prometheusAdapter:  
  operator:
    ...
    tlsEnabled: true
    tlsConfig:
      caSecret:
        name: client-ssl-certificate-1
        key: trusted-ca.crt
      certSecret:
        name: client--ssl-certificate-2
        key: client-cert.crt
      keySecret:
        name: client-ssl-certificate-3
        key: client-key.key
...
```

### Example for victoriametrics to use cert-manager

```yaml
...
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: "dev-cluster-issuer"
  vmAgent:
    ...
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmagent-tls-secret"
  vmAlert:
    ...
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmalert-tls-secret"
  vmAlertManager:
    ...
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmalertmanager-tls-secret"
  vmSingle:
    ...
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmsingle-tls-secret"
  vmAuth:
    ...
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmauth-tls-secret"
...
```

### Example for victoriametrics to use createSecret

It is important to disable automatic certificate generation from cert-manager when using `createSecret`.
`tlsConfig.generateCerts.enabled` must be set to `false` explicitly.

```yaml
...
victoriametrics:
  tlsEnabled: true
  vmAgent:
    ...
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        cert: |-
          -----BEGIN RSA PRIVATE KEY-----
          ...
          -----END RSA PRIVATE KEY-----
        secretName: "vmagent-tls-secret"
  vmAlert:
    ...
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        cert: |-
          -----BEGIN RSA PRIVATE KEY-----
          ...
          -----END RSA PRIVATE KEY-----
        secretName: "vmalert-tls-secret"
  vmAlertManager:
    ...
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        cert: |-
          -----BEGIN RSA PRIVATE KEY-----
          ...
          -----END RSA PRIVATE KEY-----
        secretName: "vmalertmanager-tls-secret"
  vmSingle:
    ...
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        cert: |-
          -----BEGIN RSA PRIVATE KEY-----
          ...
          -----END RSA PRIVATE KEY-----
        secretName: "vmsingle-tls-secret"
  vmAuth:
    ...
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        cert: |-
          -----BEGIN RSA PRIVATE KEY-----
          ...
          -----END RSA PRIVATE KEY-----
        secretName: "vmauth-tls-secret"
...
```

### Example for victoriametrics to use existingSecret

If you already have secrets created in the server containing valid certificates, `existingSecret` section can be used.

```yaml
...
victoriametrics:
  tlsEnabled: true
  vmAgent:
    ...
    tlsConfig:      
      existingSecret: "vmagent-tls-secret"
  vmAlert:
    ...
    tlsConfig:
      existingSecret: "vmalert-tls-secret"
  vmAlertManager:
    ...
    tlsConfig:      
      existingSecret: "vmalertmanager-tls-secret"
  vmSingle:
    ...
    tlsConfig:      
      existingSecret: "vmsingle-tls-secret"
  vmAuth:
    ...
    tlsConfig:      
      existingSecret: "vmauth-tls-secret"
...
```


[Back to TOC](#table-of-content)

