# Legacy Grafana dashboard

**[Back](../../README.md)**

This example shows a legacy `integreatly.org/v1alpha1` `GrafanaDashboard` and the
`grafana.integreatly.org/v1beta1` dashboard that the [Grafana converter][converter] creates from it. Use it to check
that legacy dashboards still reach Grafana after the upgrade to Grafana Operator v5.

New dashboards should use the v1beta1 API directly, as in the
[Grafana dashboard example](../grafana-dashboard/README.md).

## Files

* [Legacy dashboard](legacy-dashboard.yaml)

## Apply the example

Apply the legacy dashboard to a namespace that the converter watches:

```bash
kubectl apply -n <namespace> -f legacy-dashboard.yaml
```

## Check the converted dashboard

The converter creates a v1beta1 dashboard with the same name in the same namespace:

```bash
kubectl get grafanadashboards.grafana.integreatly.org legacy-simple-dashboard -n <namespace> -o yaml
```

The converted dashboard looks like this, with the dashboard JSON shortened:

```yaml
apiVersion: grafana.integreatly.org/v1beta1
kind: GrafanaDashboard
metadata:
  name: legacy-simple-dashboard
  namespace: <namespace>
  labels:
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by-operator: grafana-operator-converter
spec:
  allowCrossNamespaceImport: true
  folder: Legacy dashboards
  instanceSelector:
    matchLabels:
      app.kubernetes.io/component: grafana
      app.kubernetes.io/part-of: monitoring
  json: |
    { "uid": "legacy-simple-dashboard", "title": "Legacy Simple Dashboard", ... }
status:
  conditions:
    - type: DashboardSynchronized
      status: "True"
```

`DashboardSynchronized` is `True` when Grafana has loaded the dashboard. Open Grafana and find
`Legacy Simple Dashboard` in the `Legacy dashboards` folder.

## Update and delete

* Edit the legacy dashboard, and the converter updates the converted dashboard.
* Delete the legacy dashboard, and the converted dashboard stays. Delete it manually, or enable
  `grafanaConverter.grafana.converter.deleteTargetOnSourceDeletion` to delete it together with the legacy dashboard.

```bash
kubectl delete -n <namespace> -f legacy-dashboard.yaml
kubectl delete grafanadashboards.grafana.integreatly.org legacy-simple-dashboard -n <namespace>
```

## Links

* [Grafana converter](../../../installation/components/grafana-stack/grafana-converter.md)
* [qubership-grafana-operator-converter](https://github.com/Netcracker/qubership-grafana-operator-converter)
* [Grafana Operator dashboard API](https://grafana.github.io/grafana-operator/docs/api/#grafanadashboardspec)

[converter]: ../../../installation/components/grafana-stack/grafana-converter.md
