# Grafana converter deploy parameters

Helm values for the [Grafana converter](../../../installation/components/grafana-stack/grafana-converter.md).

| File                                                       | Use it when                                                                                   |
| ---------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [privileged-values.yaml](privileged-values.yaml)           | The release has cluster-wide rights. Converts selected legacy dashboards from all namespaces. |
| [non-privileged-values.yaml](non-privileged-values.yaml)   | The release has no cluster-wide rights. Converts legacy dashboards from listed namespaces.    |
| [disabled-values.yaml](disabled-values.yaml)               | No product publishes legacy dashboards.                                                       |

Apply a file on top of your other values:

```bash
helm upgrade --install monitoring-operator charts/qubership-monitoring-operator \
  --namespace monitoring \
  --skip-crds \
  --values docs/examples/deploy-parameters/grafana-converter/privileged-values.yaml
```

For a namespaced guest beside a privileged host, use the
[namespaced guest example](../namespaced-guest/README.md) instead.
