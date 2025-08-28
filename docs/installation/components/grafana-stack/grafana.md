### grafana

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                            | Scheme                                                                                                                                                |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| install                    | Allows to disable deploy Grafana. If Grafana was not deployed during the deployment using helm, it can be deployed using the change custom resource PlatformMonitoring.                                                | bool                                                                                                                                                  |
| paused                     | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                                                  |
| image                      | A docker image to be used for the grafana deployment.                                                                                                                                                                  | string                                                                                                                                                |
| ingress                    | Allows to create Ingress for Grafana UI using monitoring-operator.                                                                                                                                                     | [v1.Ingress](#ingress)                                                                                                                                |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)                          |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)                             |
| dataStorage                | Allows set a means to configure the grafana data storage.                                                                                                                                                              | [grafv1alpha1.GrafanaDataStorage](https://github.com/grafana/grafana-operator/blob/v4/documentation/deploy_grafana.md#configuring-data-storage) |
| extraVars                  | Allows set extra system environment variables for grafana.                                                                                                                                                             | map[string]string                                                                                                                                     |
| grafanaHomeDashboard       | Allows set custom home dashboard for grafana. Dependence: `extraVars: GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH: /etc/grafana-configmaps/grafana-home-dashboard/grafana-home-dashboard.json`                           | bool                                                                                                                                                  |
| backupDaemonDashboard      | Enables Backup Daemon Dashboard installation.                                                                                                                                                                          | bool                                                                                                                                                  |
| dashboardLabelSelector     | Allows to query over a set of resources according to labels.<br/>The result of matchLabels and matchExpressions are ANDed.<br/>An empty label selector matches all objects. A null label selector matches no objects.  | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)                                 |
| dashboardNamespaceSelector | Allows to query over a set of resources in namespaces that fits label selector.                                                                                                                                        | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)                                   |
| podMonitor                 | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                                                  |
| config                     | Allows set configuration for grafana. The properties used to generate grafana.ini.                                                                                                                                     | [grafv1alpha1.GrafanaConfig](https://github.com/grafana/grafana-operator/blob/v4/documentation/deploy_grafana.md#config-reconciliation)  |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                                                     |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                                                     |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                   | string                                                                                                                                                |
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


#### grafana-operator

<!-- markdownlint-disable line-length -->
| Field              | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| image              | A docker image to be used for the grafana-operator deployment.                                                                                                                                                         | string                                                                                                                       |
| paused             | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                         |
| initContainerImage | A docker image to be used into initContainer in the Grafana deployment.                                                                                                                                                | string                                                                                                                       |
| resources          | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext    | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| podMonitor         | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                         |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations        | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels             | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName  | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                   | string                                                                                                                       |
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


#### grafana-image-renderer

<!-- markdownlint-disable line-length -->
**Warning**: The grafana-image-renderer requires two extra environment variables in Grafana:

* GF_RENDERING_SERVER_URL - `http://<image-renderer-address>:<port>/render`
* GF_RENDERING_CALLBACK_URL - `http://<grafana-adderss>:<port>/`

These variables have been set by default for local renderer and Grafana services. You don't have to override them. You
need change them in case if youare yousing external renderer.

**Warning**: Rendering images requires a lot of memory, mainly because Grafana creates browser instances in the
background for the actual rendering. If you are going to render a lot of panels it make sense allocate much more memory
than default value(developers of plugin suggest 16GB ram).

| Field             | Description                                                                                                                                                                                                                                                                     | Scheme                                                                                                                       |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install           | Allows to enable deploy Grafana image renderer.                                                                                                                                                                                                                                 | *bool                                                                                                                        |
| image             | A docker image to use for grafana-image-renderer deployment.                                                                                                                                                                                                                    | string                                                                                                                       |
| name              | This name is used as the name of the microservice deployment and in labels.                                                                                                                                                                                                     | []string                                                                                                                     |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                              | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                          | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                | map[string]string                                                                                                            |
| resources         | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                                                                           | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| tolerations       | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                         | []v1.Toleration                                                                                                              |
| port              | Port for grafana-image-renderer service.                                                                                                                                                                                                                                        | integer                                                                                                                      |
| extraEnvs         | Allow to set extra system environment variables for grafana-image-renderer. More information about env  variables in [Configuration guide](https://grafana.com/docs/grafana/v9.0/setup-grafana/image-rendering/?src=your_stories_page---------------------------#configuration) | map[string]string                                                                                                            |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                              | map[string]string                                                                                                            |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                           | string                                                                                                                       |
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


