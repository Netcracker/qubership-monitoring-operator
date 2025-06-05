This document describes which abilities monitoring has to collect metrics from services/microservices with endpoints
closed by various auth types.

# Before you begin

* All credentials and tokens **must** store in Kubernetes Secrets
* All Secrets **must** be created in the namespace where was create `ServiceMonitor`/`PodMonitor`
* Auth settings **must** be configured per **each endpoint** in `ServiceMonitor`/`PodMonitor` that requires to auth

# General recommendation

Although Monitoring supports metrics collection from endpoint secured by various auth type, it doesn't means that
you need to close all your metrics endpoint by auth.

We almost all recommended following two rules for metrics and metrics endpoint:

* Metrics **shouldn't** contain any sensitive data that must require to be secured
* Endpoint with metrics **shouldn't** require any authentication

But if you for some reason need to secure the endpoint with metrics we recommend using RBAC auth. And particular we
recommend using authentication by `nonResourceURL`.

Monitoring has special Cluster permissions on:

```yaml
- nonResourceURLs:
    - /metrics
  verbs:
    - 'get'
```

So you can configure your service to authenticate Prometheus or VMAgent who will try to collect metrics from the
service with such permissions.

You can use the third-party `kube-rbac-proxy` if you don't want to add such logic to your application.
GitHub:

* [https://github.com/brancz/kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy)

And install it as a sidecar to your application.

# Authentication methods

The application can support different auth methods, like BasicAuth, OAuth2 and so on. And deploy this application and
add some additional auth methods (that will work only in Kubernetes) like using auth by ServiceAccount (SA) permissions.

So we can specify the following popular auth types:

* BasicAuth
* Static token
* OAuth2
* Custom URL params
* RBAC auth

So let's describe how to configure Monitoring to collect metrics from the endpoint that uses each of these auth types.

## Basic Auth

The Custom Resources `ServiceMonitor` and `PodMonitor` allow referring to the Secret with BasicAuth credentials
and use them to pull metrics.

To configure it you need to use a section `basicAuth` for each endpoint in CRs.

Official `prometheus-operator` documentation:

* [ServiceMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#endpoint)
* [PodMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#monitoring.coreos.com/v1.PodMetricsEndpoint)

For example, suppose we have a Secret:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: test-service-credentials
type: Opaque
data:
  username: YWRtaW4=   # "admin" in base64
  password: YWRtaW4=   # "admin" in base64
```

So to collect metrics from the endpoint you can configure `ServiceMonitor` or `PodMonitor` as follows:

`ServiceMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that BasicAuth config should be specify for each endpoint if it use BasicAuth
    - interval: 30s
      path: /metrics
      basicAuth:     # <-- Need to add this section
        username:
          name: test-service-credentials
          key: username
        password:
          name: test-service-credentials
          key: password
```

`PodMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that BasicAuth config should be specify for each endpoint if it use BasicAuth
    - interval: 30s
      path: /metrics
      basicAuth:     # <-- Need to add this section
        username:
          name: test-service-credentials
          key: username
        password:
          name: test-service-credentials
          key: password
```

## Static token

The Custom Resources `ServiceMonitor` and `PodMonitor` allow referring to the Secret using static tokens.

To configure it you need to use a section `bearerTokenSecret` or `authorization` for each endpoint in CRs.

Official `prometheus-operator` documentation:

* [ServiceMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#endpoint)
* [PodMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#monitoring.coreos.com/v1.PodMetricsEndpoint)
* [Authorization](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#monitoring.coreos.com/v1.SafeAuthorization)

For example, suppose we have a Secret:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: test-service-token
type: Opaque
data:
  token: bGprZ2h1aUdieXVpaGpnZ2hqZmRBU0RBS0xI   # "ljkghuiGbyuihjgghjfdASDAKLH" in base64
```

So to collect metrics from the endpoint you can configure `ServiceMonitor` or `PodMonitor` as follows:

`ServiceMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      bearerTokenSecret:     # <-- Need to add this section
        name: test-service-token
        key: token
```

or

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      authorization:     # <-- Need to add this section
        type: Bearer
        credentials:
          name: test-service-token
          key: token
```

`PodMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      bearerTokenSecret:     # <-- Need to add this section
        name: test-service-token
        key: token
```

or

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      authorization:     # <-- Need to add this section
        type: Bearer
        credentials:
          name: test-service-token
          key: token
```

## OAuth2

The Custom Resources `ServiceMonitor` and `PodMonitor` allow referring to the Secret using static tokens.

To configure it you need to use a section `oauth2` for each endpoint in CRs.

Official `prometheus-operator` documentation:

* [ServiceMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#endpoint)
* [PodMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#monitoring.coreos.com/v1.PodMetricsEndpoint)
* [OAuth2](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#monitoring.coreos.com/v1.OAuth2)

For example, suppose we have a Secret:

**WARNING!** In this example, we will use integration and link to Keycloak.

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: test-service-oauth2-creds
type: Opaque
data:
  clientId: dGVzdA==                                           # "test" in base64
  clientSecret: ZE1EdU9kVHJpQWRiOTFxcGFveFdjeUlsQzFPYklHM1k=   # "dMDuOdTriAdb91qpaoxWcyIlC1ObIG3Y" in base64
```

So to collect metrics from the endpoint you can configure `ServiceMonitor` or `PodMonitor` as follows:

`ServiceMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      oauth2:     # <-- Need to add this section
        clientId:
          secret:
            name: test-service-oauth2-creds
            key: clientId
        clientSecret:
          name: test-service-oauth2-creds
          key: clientSecret
        tokenUrl: http://<keycloak_url>/realms/envoy/protocol/openid-connect/token
        scopes:
          - user
```

`PodMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      oauth2:     # <-- Need to add this section
        clientId:
          secret:
            name: test-service-oauth2-creds
            key: clientId
        clientSecret:
          name: test-service-oauth2-creds
          key: clientSecret
        tokenUrl: http://<keycloak_url>/realms/envoy/protocol/openid-connect/token
        scopes:
          - user
```

## Custom URL params

**NOTE:** Currently Prometheus has no ability to scrape metrics from endpoints with custom headers. Only with using
custom URL parameters.

In case your service requires some custom URL parameters to auth users you can specify them in CRs like
`ServiceMonitor` or `PodMonitor`.

Official `prometheus-operator` documentation:

* [ServiceMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#endpoint)
* [PodMonitor: Endpoint](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#monitoring.coreos.com/v1.PodMetricsEndpoint)

So to collect metrics from the endpoint you can configure `ServiceMonitor` or `PodMonitor` as follows:

`ServiceMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      params:     # <-- Need to add this section
        custom-user: admin
        page: 123
```

`PodMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
...
spec:
  endpoints:
    # Pay attention, that token config should be specify for each endpoint if it use token
    - interval: 30s
      path: /metrics
      params:     # <-- Need to add this section
        custom-user: admin
        page: 123
```

## RBAC auth

As we already mentioned in the [General recommendation](#general-recommendation) section Monitoring has special
permissions that it can use for authenticating using RBAC permissions (based on permissions from ServiceAccount).

If the microservice will this auth type it can check any permissions from Role or ClusterRole
(bind with ServiceAccount) of the metric collector. But we advise using `nonResourceURL` for the URL:

```yaml
- nonResourceURLs:
    - /metrics
  verbs:
    - 'get'
```

This approach requires configuring auth check on the microservice side but doesn't require any settings
on the Monitoring side or in Custom Resources.
