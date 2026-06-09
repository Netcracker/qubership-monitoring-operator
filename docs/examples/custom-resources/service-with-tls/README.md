# ServiceMonitor with using https and TLS

* [ServiceMonitor with using https and TLS](#servicemonitor-with-using-https-and-tls)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show a basic configuration for ServiceMonitor.

```yaml
spec:
  endpoints:
  - interval: 30s
    port: metrics
    scheme: https
    tlsConfig:
      caFile: /etc/prometheus/secrets/kube-etcd-client-certs/etcd-client-ca.crt
      certFile: /etc/prometheus/secrets/kube-etcd-client-certs/etcd-client.crt
      keyFile: /etc/prometheus/secrets/kube-etcd-client-certs/etcd-client.key
      serverName: ""
  jobLabel: k8s-app
  namespaceSelector:
    matchNames:
    - kube-system
  selector:
    matchLabels:
      k8s-app: etcd
```

It means that Prometheus will collect metrics from service with settings:

* metrics will collect with job with label `k8s-app`
* metrics will collect from all pods with label `k8s-app: etcd`
* metrics will collect from all discovered pod from port with name `metrics` with interval `30s` and schema `https`
* metrics will collect with TLS config and specified certificates
* for verify certificates will use name "" (specify in `serverName: ""`)

## Files

* [Service Monitor Etcd](service-monitor-etcd.yaml)
* [Service Monitor Kubelet](service-monitor-kubelet.yaml)

## How to apply example

Kubernetes:

```bash
kubectl apply -f service-monitor-etcd.yaml
```

OpenShift:

```bash
oc apply -f service-monitor-etcd.yaml
```

## Links

Prometheus operator API

* [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md)

Victoriametrics operator API

* [API Documentation](https://docs.victoriametrics.com/operator/api.html)
