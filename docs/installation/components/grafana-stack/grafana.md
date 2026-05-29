# Grafana

Sensitive parameters are configured via files mounted from Kubernetes Secrets:

- `GF_SECURITY_ADMIN_USER` and `GF_SECURITY_ADMIN_PASSWORD` are read from `/{grafana-name}-admin-credentials` Secret using Grafana file provider (`$__file{...}`).
- `GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET` is read from `grafana-oauth-client-secret` Secret using Grafana file provider (`$__file{...}`).
- `GF_AUTH_GENERIC_OAUTH_CLIENT_ID` remains a non-sensitive value and may stay in environment variables.

## Generic OAuth (IdP integration)

When `spec.auth` is populated in `PlatformMonitoring`, the Monitoring Operator configures Grafana's
built-in [Generic OAuth](https://grafana.com/docs/grafana/latest/setup-grafana/configure-access/configure-authentication/generic-oauth/)
provider by propagating the following fields into `spec.config["auth.generic_oauth"]` of the Grafana CR:

| `spec.auth` field         | Grafana ini key               | Notes                                                            |
|---------------------------|-------------------------------|------------------------------------------------------------------|
| `loginUrl`                | `auth_url`                    | IdP authorization endpoint                                       |
| `tokenUrl`                | `token_url`                   | IdP token endpoint                                               |
| `userInfoUrl`             | `api_url`                     | IdP userinfo endpoint                                            |
| _(always set)_            | `enabled` = `"true"`          | Activates Generic OAuth                                          |
| _(always set)_            | `scopes` = `"openid profile"` | Required OIDC scopes                                             |
| `clientSecret` (via Helm) | `client_secret`               | Read from `grafana-oauth-client-secret` Secret via file provider |

### TLS configuration for OAuth

When `spec.auth.tlsConfig` is set, the Monitoring Operator mounts the referenced Kubernetes Secrets
as read-only volumes in the Grafana pod and configures the corresponding Grafana ini keys:

| `spec.auth.tlsConfig` field | Volume mount path                      | Grafana ini key            |
|-----------------------------|----------------------------------------|----------------------------|
| `caSecret`                  | `/etc/grafana-tls/<secret-name>/<key>` | `tls_client_ca`            |
| `certSecret`                | `/etc/grafana-tls/<secret-name>/<key>` | `tls_client_cert`          |
| `keySecret`                 | `/etc/grafana-tls/<secret-name>/<key>` | `tls_client_key`           |
| `insecureSkipVerify`        | _(no mount)_                           | `tls_skip_verify_insecure` |

The referenced Secrets must exist in the same namespace as the `PlatformMonitoring` resource.
Each unique Secret name results in one Volume mounted at `/etc/grafana-tls/<secret-name>/`.
If `certSecret` and `keySecret` reference the same Secret name, it is mounted only once.

Example `spec.auth` configuration:

```yaml
auth:
  loginUrl: https://keycloak.example.com/realms/master/protocol/openid-connect/auth
  tokenUrl: https://keycloak.example.com/realms/master/protocol/openid-connect/token
  userInfoUrl: https://keycloak.example.com/realms/master/protocol/openid-connect/userinfo
  tlsConfig:
    insecureSkipVerify: false
    caSecret:
      name: idp-ca-cert
      key: ca.crt
    certSecret:
      name: idp-client-cert
      key: tls.crt
    keySecret:
      name: idp-client-cert
      key: tls.key
```

## Admin credentials

| Topic          | Details                                                                                                                                                          |
|----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Secret name    | `{grafana-name}-admin-credentials` (default: `grafana-admin-credentials`)                                                                                        |
| Created by     | Helm when `grafana.disableDefaultAdminSecret=false` (default), from `grafana.security.admin_user` / `admin_password` in chart values (default `admin` / `admin`) |
| Consumed by    | Grafana via mounted files; Monitoring Operator syncs runtime password when the Secret changes                                                                    |
| Runtime change | Edit the Secret; see [password change guide](../../../user-guides/password-change.md)                                                                            |

On **first install**, Grafana applies credentials from the mounted Secret when the admin user is created in the database.
When the Secret is **updated later**, Monitoring Operator sets annotation `checksum/admin-secret` on the Grafana CR
pod template (rolling restart) and runs `grafana cli admin reset-admin-password` so the password in the database
matches the Secret (required when persistent storage is enabled).

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                              | Scheme                                                                                                                                          |
|----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| install                    | Allows to disable deploy Grafana. If Grafana was not deployed during the deployment using helm, it can be deployed using the change custom resource PlatformMonitoring.                                                  | bool                                                                                                                                            |
| paused                     | Set paused to reconciliation.                                                                                                                                                                                            | bool                                                                                                                                            |
| image                      | A Docker image to be used for the grafana deployment.                                                                                                                                                                    | string                                                                                                                                          |
| ingress                    | Allows to create Ingress for Grafana UI using monitoring-operator.                                                                                                                                                       | *[Ingress](../../../api/platform-monitoring.md#ingress)                                                                                         |
| httpRoute                  | HTTPRoute allows to create Gateway API HTTPRoute for the Grafana UI.                                                                                                                                                     | [HTTPRouteSpec](https://gateway-api.sigs.k8s.io/reference/api-spec/main/spec/#httproute)                                                        |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                    | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)                    |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                       | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)                       |
| dataStorage                | Allows set a means to configure the grafana data storage.                                                                                                                                                                | [grafv1alpha1.GrafanaDataStorage](https://github.com/grafana/grafana-operator/blob/v4/documentation/deploy_grafana.md#configuring-data-storage) |
| extraVars                  | Allows set extra system environment variables for grafana.                                                                                                                                                               | map[string]string                                                                                                                               |
| grafanaHomeDashboard       | Allows set custom home dashboard for grafana. Dependence: `extraVars: GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH: /etc/grafana-configmaps/grafana-home-dashboard/grafana-home-dashboard.json`                             | bool                                                                                                                                            |
| backupDaemonDashboard      | Enables Backup Daemon Dashboard installation.                                                                                                                                                                            | bool                                                                                                                                            |
| dashboardLabelSelector     | Allows to query over a set of resources according to labels.<br/>The result of matchLabels and matchExpressions are ANDed.<br/>An empty label selector matches all objects. A null label selector matches no objects.    | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)                           |
| dashboardNamespaceSelector | Allows to query over a set of resources in namespaces that fits label selector.                                                                                                                                          | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)                             |
| podMonitor                 | Pod monitor for self monitoring.                                                                                                                                                                                         | *[Monitor](../../../api/platform-monitoring.md#monitor)                                                                                         |
| config                     | Allows set configuration for grafana. The properties used to generate grafana.ini.                                                                                                                                       | [grafv1alpha1.GrafanaConfig](https://github.com/grafana/grafana-operator/blob/v4/documentation/deploy_grafana.md#config-reconciliation)         |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                           | *v1.Affinity                                                                                                                                    |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"   | map[string]string                                                                                                                               |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                         | map[string]string                                                                                                                               |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                     | string                                                                                                                                          |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafana:
  install: true
  paused: false
  image: grafana/grafana:11.6.5
  ingress:
    ...
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  config:
    auth:
      disable_login_form: false
      disable_signout_menu: true
    auth.anonymous:
      enabled: false
    log:
      level: warn
      mode: console
  extraVars:
    GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH: /etc/grafana-configmaps/grafana-home-dashboard/grafana-home-dashboard.json
    GF_LIVE_ALLOWED_ORIGINS: "*"
    GF_FEATURE_TOGGLES_ENABLE: ngalert
  grafanaHomeDashboard: true
  backupDaemonDashboard: true
  dashboardLabelSelector:
    - matchLabels:
        app.kubernetes.io/component: monitoring
      matchExpressions:
        - key: openshift.io/cluster-monitoring
          operator: NotIn
          values: [ "true" ]
    - matchExpressions:
        - key: app.kubernetes.io/instance
          operator: Exists
        - key: app.kubernetes.io/version
          operator: Exists
  dashboardNamespaceSelector:
    matchLabels:
      label-key: label-value
    matchExpressions:
      - key: openshift.io/cluster-monitoring
        operator: NotIn
        values: [ "true" ]
      - key: kubernetes.io/metadata.name
        operator: In
        values:
          - monitoring
          - cassandra
  podMonitor:
    ...see example by link...
  dataStorage:
    labels:
      app: grafana
    annotations:
      app: grafana
    accessModes:
      - ReadWriteOnce
    size: 2Gi
    class: local-storage
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```


## grafana-operator

<!-- markdownlint-disable line-length -->
| Field              | Description                                                                                                                                                                                                            | Scheme                                                                                                                        |
|--------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| image              | A Docker image to be used for the grafana-operator deployment.                                                                                                                                                         | string                                                                                                                        |
| paused             | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                          |
| initContainerImage | A Docker image to be used into initContainer in the Grafana deployment.                                                                                                                                                | string                                                                                                                        |
| resources          | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)  |
| securityContext    | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)     |
| podMonitor         | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](../../../api/platform-monitoring.md#monitor)                                                                       |
| affinity           | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                  |
| annotations        | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                             |
| labels             | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                             |
| priorityClassName  | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                   | string                                                                                                                        |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafana:
  operator:
    image: integreatly/grafana-operator:latest
    paused: false
    initContainerImage: integreatly/grafana_plugins_init:latest
    resources:
      limits:
        cpu: 200m
        memory: 200Mi
      requests:
        cpu: 100m
        memory: 100Mi
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    podMonitor:
      ...see example by link...
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
```


## grafana-image-renderer

<!-- markdownlint-disable line-length -->
**Warning**: The grafana-image-renderer requires two extra environment variables in Grafana:

* GF_RENDERING_SERVER_URL - `http://<image-renderer-address>:<port>/render`
* GF_RENDERING_CALLBACK_URL - `http://<grafana-adderss>:<port>/`

These variables have been set by default for local renderer and Grafana services. You don't have to override them. You
need change them in case if youare yousing external renderer.

**Warning**: Rendering images requires a lot of memory, mainly because Grafana creates browser instances in the
background for the actual rendering. If you are going to render a lot of panels it make sense allocate much more memory
than default value(developers of plugin suggest 16GB ram).

| Field             | Description                                                                                                                                                                                                                                                                      | Scheme                                                                                                                       |
|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------|
| install           | Allows to enable deploy Grafana image renderer.                                                                                                                                                                                                                                  | *bool                                                                                                                        |
| image             | A Docker image to use for grafana-image-renderer deployment.                                                                                                                                                                                                                     | string                                                                                                                       |
| name              | This name is used as the name of the microservice deployment and in labels.                                                                                                                                                                                                      | []string                                                                                                                     |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                               | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                           | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                 | map[string]string                                                                                                            |
| resources         | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                                                                            | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| tolerations       | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                          | []v1.Toleration                                                                                                              |
| port              | Port for grafana-image-renderer service.                                                                                                                                                                                                                                         | integer                                                                                                                      |
| extraEnvs         | Allow to set extra system environment variables for grafana-image-renderer. More information about env  variables in [Configuration guide](https://grafana.com/docs/grafana/v9.0/setup-grafana/image-rendering/?src=your_stories_page---------------------------#configuration)  | map[string]string                                                                                                            |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                               | map[string]string                                                                                                            |
| affinity          | If specified, the pod's scheduling constraints                                                                                                                                                                                                                                   | *v1.Affinity                                                                                                                 |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                            | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafana:
  imageRenderer:
    install: true
    image: grafana/grafana-image-renderer:3.12.9
    name: grafana-image-renderer
    resources:
      limits:
        cpu: 300m
        memory: 500Mi
      requests:
        cpu: 150m
        memory: 250Mi
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    port: 8282
    priorityClassName: priority-class
```

