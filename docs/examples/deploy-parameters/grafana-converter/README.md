# Grafana converter deploy parameters

Helm values for the [Grafana converter](../../../installation/components/grafana-stack/grafana-converter.md).

| File                                                       | Use it when                                                                                   |
| ---------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [privileged-values.yaml](privileged-values.yaml)           | The release has cluster-wide rights. Converts selected legacy dashboards from all namespaces. |
| [non-privileged-values.yaml](non-privileged-values.yaml)   | The release has no cluster-wide rights. Converts legacy dashboards from listed namespaces.    |
| [disabled-values.yaml](disabled-values.yaml)               | You need no built-in dashboards, and no product publishes legacy dashboards.                  |

The built-in dashboards also reach Grafana through the converter. The operator and the chart publish them as legacy
dashboards in the release namespace. With the converter disabled, or with the release namespace outside its watch
scope or its dashboard selectors, Grafana has no built-in dashboards.

If you run with `grafanaConverter.install=false`, enable the converter before you upgrade, and check that its watch
scope and selectors include the release namespace.

Apply a file on top of your other values:

```bash
helm upgrade --install monitoring-operator charts/qubership-monitoring-operator \
  --namespace monitoring \
  --skip-crds \
  --values docs/examples/deploy-parameters/grafana-converter/privileged-values.yaml
```

For a namespaced guest beside a privileged host, use the
[namespaced guest example](../namespaced-guest/README.md) instead.
