This document describes how to configure authentication (and available options) for all UIs available
in the Monitoring stack.

# List of supported auth methods

Monitoring Stack contains various third-party tools with various UIs. And these UIs support different auth methods.
And not always one UIs supports the same authentication method as the other UI.

* Use Ingress annotation to add auth by Basic Auth and even via the OAuth2 provider
* Use built-in auth abilities (like auth in Grafana or auth-proxy for other components)

## Ingress annotations

Almost all Ingress controllers allow adding authentication using annotations in Ingress objects.

But in almost all our clouds we are using Ingress-Nginx, so all next configurations and steps will be based on its
functions and will refer to its documentation.

Ingress-Nginx supporta the following auth options:

* [Basic Authentication](https://kubernetes.github.io/ingress-nginx/examples/auth/basic/)
* [Client Certificate Authentication](https://kubernetes.github.io/ingress-nginx/examples/auth/client-certs/)
* [External Basic Authentication](https://kubernetes.github.io/ingress-nginx/examples/auth/external-auth/)
* [External OAuth Authentication](https://kubernetes.github.io/ingress-nginx/examples/auth/oauth-external-auth/)

In this document, we will describe how to configure only Ingress BasicAuth.

### Before you begin

* All Secrets that will be usied for Ingress configuration you must create manually and before deploy Monitoring
* In case of using external providers (for BasicAuth or for OAuth) this external provider must be already deployed
  and running before deploy Monitoring
* In case of using external providers (for BasicAuth or for OAuth), users must be created in these providers

### Restrictions

* Ingress allows adding only authentication feature
* Authorization doesn't support using Ingress auth
* Services and other links inside the cloud to Monitoring tools will work without auth
* Other Ingress controllers may not support all these features

### Ingress BasicAuth

This option allows using of authentication credentials saved in `.htpasswd` format in the Secret.

To save/encode credentials in `.htpasswd` you need to use the following command:

```bash
$ htpasswd -c auth <username>
New password: <password>
...
```

for example, for use with name `foo` and password `bar` it will look like as:

```bash
$ htpasswd -c auth foo
New password: <input "bar">
New password: <input "bar">
Re-type new password:
Adding password for user foo
```

This command will create a file with the name `auth`. To create a Secret with content from `auth` file you can use
the command:

```bash
kubectl create secret generic monitoring-basic-auth --from-file=auth --namespace=<monitoring_namespace>
```

**WARNING!** Secret **must** be created in the namespace with monitoring where will create Ingresses.

Next, you need to add the following parameters in the Monitoring deploy parameters.

#### Grafana

```yaml
grafana:
  ingress:
    install: true
    host: <you_host>
    annotations:
      # type of authentication
      nginx.ingress.kubernetes.io/auth-type: basic
      # name of the secret that contains the user/password definitions
      nginx.ingress.kubernetes.io/auth-secret: <secret_name>
      # message to display with an appropriate context why the authentication is required
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
grafana:
  ingress:
    install: true
    host: <you_host>
    annotations:
      nginx.ingress.kubernetes.io/auth-type: basic
      nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

#### Prometheus Stack

This section describes how to configure BasicAuth using Ingress for components from the Prometheus stack.

##### Alertmanager

```yaml
alertManager:
  ingress:
    install: true
    host: <you_host>
    annotations:
      # type of authentication
      nginx.ingress.kubernetes.io/auth-type: basic
      # name of the secret that contains the user/password definitions
      nginx.ingress.kubernetes.io/auth-secret: <secret_name>
      # message to display with an appropriate context why the authentication is required
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
alertManager:
  ingress:
    install: true
    host: <you_host>
    annotations:
      nginx.ingress.kubernetes.io/auth-type: basic
      nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

##### Prometheus

```yaml
prometheus:
  ingress:
    install: true
    host: <you_host>
    annotations:
      # type of authentication
      nginx.ingress.kubernetes.io/auth-type: basic
      # name of the secret that contains the user/password definitions
      nginx.ingress.kubernetes.io/auth-secret: <secret_name>
      # message to display with an appropriate context why the authentication is required
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
prometheus:
  ingress:
    install: true
    host: <you_host>
    annotations:
      nginx.ingress.kubernetes.io/auth-type: basic
      nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

#### VictoriaMetrics Stack

This section describes how to configure BasicAuth using Ingress for components from the VictoriaMetrics stack.

##### VMAgent

```yaml
victoriametrics:
  vmAgent:
    ingress:
      install: true
      host: <you_host>
      annotations:
        # type of authentication
        nginx.ingress.kubernetes.io/auth-type: basic
        # name of the secret that contains the user/password definitions
        nginx.ingress.kubernetes.io/auth-secret: <secret_name>
        # message to display with an appropriate context why the authentication is required
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
victoriametrics:
  vmAgent:
  ingress:
    install: true
    host: <you_host>
    annotations:
      nginx.ingress.kubernetes.io/auth-type: basic
      nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
      nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

##### VMAlert

```yaml
victoriametrics:
  vmAlert:
    ingress:
      install: true
      host: <you_host>
      annotations:
        # type of authentication
        nginx.ingress.kubernetes.io/auth-type: basic
        # name of the secret that contains the user/password definitions
        nginx.ingress.kubernetes.io/auth-secret: <secret_name>
        # message to display with an appropriate context why the authentication is required
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
victoriametrics:
  vmAlert:
    ingress:
      install: true
      host: <you_host>
      annotations:
        nginx.ingress.kubernetes.io/auth-type: basic
        nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

##### VMAuth

```yaml
victoriametrics:
  vmAuth:
    ingress:
      install: true
      host: <you_host>
      annotations:
        # type of authentication
        nginx.ingress.kubernetes.io/auth-type: basic
        # name of the secret that contains the user/password definitions
        nginx.ingress.kubernetes.io/auth-secret: <secret_name>
        # message to display with an appropriate context why the authentication is required
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
victoriametrics:
  vmAuth:
    ingress:
      install: true
      host: <you_host>
      annotations:
        nginx.ingress.kubernetes.io/auth-type: basic
        nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

##### VMAlertmanager

```yaml
victoriametrics:
  vmAlertManager:
    ingress:
      install: true
      host: <you_host>
      annotations:
        # type of authentication
        nginx.ingress.kubernetes.io/auth-type: basic
        # name of the secret that contains the user/password definitions
        nginx.ingress.kubernetes.io/auth-secret: <secret_name>
        # message to display with an appropriate context why the authentication is required
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
victoriametrics:
  vmAlertManager:
    ingress:
      install: true
      host: <you_host>
      annotations:
        nginx.ingress.kubernetes.io/auth-type: basic
        nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

##### VMSingle

```yaml
victoriametrics:
  vmSingle:
    ingress:
      install: true
      host: <you_host>
      annotations:
        # type of authentication
        nginx.ingress.kubernetes.io/auth-type: basic
        # name of the secret that contains the user/password definitions
        nginx.ingress.kubernetes.io/auth-secret: <secret_name>
        # message to display with an appropriate context why the authentication is required
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

for example:

```yaml
victoriametrics:
  vmSingle:
    ingress:
      install: true
      host: <you_host>
      annotations:
        nginx.ingress.kubernetes.io/auth-type: basic
        nginx.ingress.kubernetes.io/auth-secret: monitoring-basic-auth
        nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required'
```

## Built-in auth

We prepared a matrix of supported auth methods per UI:

| Component      | Basic Auth | OAuth2      | LDAP                            |
| -------------- | ---------- | ----------- | ------------------------------- |
| Grafana        | ✓ Support  | ✓ Support   | ✓ Support                       |
| Alertmanager   | ✓ Support* | ✓ Support** | ⚠ Support via OAuth2 provider** |
| Prometheus     | ✓ Support* | ✓ Support** | ⚠ Support via OAuth2 provider** |
| VMAgent        | ✓ Support  | ✓ Support** | ⚠ Support via OAuth2 provider** |
| VMAlert        | ✓ Support  | ✓ Support** | ⚠ Support via OAuth2 provider** |
| VMAuth         | ✓ Support  | ✓ Support** | ⚠ Support via OAuth2 provider** |
| VMAlertmanager | ✓ Support  | ✓ Support** | ⚠ Support via OAuth2 provider** |
| VMSingle       | ✓ Support  | ✓ Support** | ⚠ Support via OAuth2 provider** |

where:

* `*` - Prometheus and Alertmanager support Basic Auth since version `Prometheus >= 2.24.0`
  and `Alertmanager >= 0.19.0`.
* `**` - Prometheus and VictoriaMetrics components doesn't support OAuth2, but we can deploy `oauth2-proxy` as sidecar.
  This container allows add OAuth2 auth method.

If we are talking about support authentication and authorization we have the following matrix:

| Component      | Authentication                  | Authorization                         |
| -------------- | ------------------------------- | ------------------------------------- |
| Grafana        | ✓ Support (Basic, LDAP, OAuth2) | ✓ Support (View, Editor, Admin roles) |
| Prometheus     | ✓ Support (OAuth2)               | ✗ Not support                         |
| Alertmanager   | ✓ Support (OAuth2)               | ✗ Not support                         |
| VMAgent        | ✓ Support (OAuth2)               | ✗ Not support                         |
| VMAlert        | ✓ Support (OAuth2)               | ✗ Not support                         |
| VMAuth         | ✓ Support (OAuth2)               | ✗ Not support                         |
| VMAlertmanager | ✓ Support (OAuth2)               | ✗ Not support                         |
| VMSingle       | ✓ Support (OAuth2)               | ✗ Not support                         |

Details about each component see below.

### Common

In Monitoring Stack include some UIs which has different auth methods. And to simplify configuration we implemented
a common `auth` section for configuring them.

If you want to enable OAuth2 auth method for all UIs you can specify all parameters:

```yaml
auth:
  clientId: xxx
  clientSecret: xxx
  loginUrl: http://1.2.3.4/authorize
  tokenUrl: http://1.2.3.4/token
  userInfoUrl: http://1.2.3.4/userinfo

  ## Settings for TLS
  tlsConfig:
    caSecret:
      name: grafana-client-certificates
      key: trusted-ca.crt
    certSecret:
      name: grafana-client-certificates
      key: client-cert.crt
    keySecret:
      name: grafana-client-certificates
      key: client-key.key
    insecureSkipVerify: false
```

Settings from `auth` section will provide to:

* `Grafana` - settings will provide into Grafana settings because Grafana already native supported Generic OAuth
* `Prometheus` - in the Prometheus pod will add a sidecar with `oauth2-proxy` to block access to the UI
* `Alertmanager` - in the Alertmanager pod will add sidecar with `oauth2-proxy` to block access to the UI

Details and parameters description see at [Installation Guide: Auth](../installation/README.md#auth).

### Grafana

Grafana supports a lot of auth methods. Some methods allow in the OSS version, but some are available only
in the Enterprise version.

OSS auth methods:

* Basic auth
* OAuth2
* LDAP

Enterprise auth methods:

* SAML
* Advanced authorization

Officinal documentation [https://grafana.com/docs/grafana/latest/auth/](https://grafana.com/docs/grafana/latest/auth/).

In Monitoring Stack we supported all auth methods available in Grafana OSS.

There are some ways to configure various auth methods for Grafana in Monitoring Stack:

* Specify parameters in `auth` section to configure OAuth2.
* Manually specify in `grafana` section. This option allows configuring BasicAuth or LDAP.

**Warning!** In case of using OAuth2 we have two prerequisites to Identity Provider:

* Client should have `email` in client scopes
* All users that wants to login in Grafana should have filled `email` user field

Examples:

* Enable OAuth for Grafana (and other UIs)

```yaml
auth:
  clientId: xxx
  clientSecret: xxx
  loginUrl: http://1.2.3.4/authorize
  tokenUrl: http://1.2.3.4/token
  userInfoUrl: http://1.2.3.4/userinfo
```

* Enable LDAP (only for Grafana)

Add settings in PlatformMonitoring CR:

```yaml
apiVersion: monitoring.qubership.org/v1alpha1
kind: PlatformMonitoring
metadata:
  name: platformmonitoring
spec:
  grafana:
    config:
      auth.ldap:
        enabled: true
        # By this path will mount secret which described below
        config_file: /etc/grafana-secrets/grafana-ldap-config/ldap.toml
```

Add `data` content in already existing Secret `grafana-ldap-config`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: grafana-ldap-config
data:
  ldap.toml: |-
    # Add content from this line in Secret
    [[servers]]
    host = "127.0.0.1"
    port = 389
    use_ssl = false
    start_tls = false
    ssl_skip_verify = false
    bind_dn = "cn=admin,dc=grafana,dc=org"
    bind_password = 'grafana'
    search_filter = "(cn=%s)"
    search_base_dns = ["dc=grafana,dc=org"]

    [servers.attributes]
    name = "givenName"
    surname = "sn"
    username = "cn"
    member_of = "memberOf"
    email =  "email"
```

**Note:** You can change Secret data using UI (ex. Kubernetes Dashboard). But also you can
edit it using CLI tools such as `kubectl` if secrets edit in UI is not available.

For example:

```bash
kubectl edit secret grafana-ldap-config -n <namespace with monitoring>
```

More information about configuring LDAP in Grafana you can read in the documentation
[Grafana: LDAP](https://grafana.com/docs/grafana/latest/auth/ldap/).

### Prometheus Stack

This section describes how to configure auth using built-in capabilities for components from the Prometheus stack.

#### Alertmanager

Alertmanager was initially designed as a cloud-native application and all authentication and authorization features
it were shifted to other tools.

But since `Alertmanager >= 0.19.0` is starting support to the BasicAuth method.

Official documentation [https://prometheus.io/docs/guides/basic-auth/](https://prometheus.io/docs/guides/basic-auth/).

But in Cloud very often use various proxy tools which add additional auth methods.
And for Alertmanager we use `oauth2-proxy`.

Documentation about `oauth2-proxy` [https://github.com/oauth2-proxy/oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy).

Currently Monitoring Stack doesn't allow configuring auth method separately. You should use `auth` section to enable it.
And currently, we don't support Basic auth for Alertmanager.

How to configure `auth` section see at:

* [Enable authentication for UIs: Common](#common)

#### Prometheus

Prometheus was initially designed as a cloud-native application and all authentication and authorization features it were
shifted to other tools.

But since `Prometheus >= 2.24.0` is starting to support the BasicAuth method.

Official documentation [https://prometheus.io/docs/guides/basic-auth/](https://prometheus.io/docs/guides/basic-auth/).

But in Cloud very often use various proxy tools which add additional auth methods.
And for Prometheus we use `oauth2-proxy`.

Documentation about `oauth2-proxy` [https://github.com/oauth2-proxy/oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy).

Currently Monitoring Stack doesn't allow configuring auth method separately. You should use `auth` section to enable it.
And currently, we don't support Basic auth for Prometheus.

How to configure `auth` section see at:

* [Enable authentication for UIs: Common](#common)

### VictoriaMetrics Stack

This section describes how to configure auth using built-in capabilities for components from the VictoriaMetrics stack.

All VictoriaMetrics components support BasicAuth, but now we don't provide an ability to configure BasicAuth,
only OAuth2 via oauth2-proxy.

#### VMAgent

To support OAuth2 using `oauth2-proxy`. It uses configuration from the common section `auth`.

How to configure `auth` section see at:

* [Enable authentication for UIs: Common](#common)

#### VMAlert

To support OAuth2 using `oauth2-proxy`. It uses configuration from the common section `auth`.

How to configure `auth` section see at:

* [Enable authentication for UIs: Common](#common)

#### VMAlertmanager

To support OAuth2 using `oauth2-proxy`. It uses configuration from the common section `auth`.

How to configure `auth` section see at:

* [Enable authentication for UIs: Common](#common)

#### VMSingle

To support OAuth2 using `oauth2-proxy`. It uses configuration from the common section `auth`.

How to configure `auth` section see at:

* [Enable authentication for UIs: Common](#common)

