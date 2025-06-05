This document describes all exporters from which the Monitoring stack can collect metrics or juts a list of metrics
available to collect from some tools.

# Types

## Applications

* Java
  * Micrometer
  * JMX exporter
* Golang
* Python
* Nginx

## Big Data

* Airflow
* Apache Spark

## Databases

* ArangoDB
* Cassandra
* ClickHouse
* ElasticSearch
* Greenplum
* MongoDB
* PostgreSQL

## Distributes Caches

* [Consul](metrics/consul-metrics.md)
* Redis
* ZooKeeper

## Kubernetes

* Calico
* [Etcd](metrics/etcd-metrics.md)
* Ingress Nginx
* Kube APIServer
* Kube Scheduler
* Kube-state-metrics
* [Kubelet](metrics/kubelet-metrics.md)

## Observability

### Components

* [Alertmanager](metrics/alertmanager-metrics.md)
* FluentBit
* FluentD
* Graylog
* [Prometheus](metrics/prometheus-metrics.md)
* Jaeger

### Exporters

* [Blackbox-exporter](exporters/blackbox-exporter.md)
* Cert-exporter
* [JSON-exporter](exporters/json-exporter.md)
* Network-latency-exporter
* Version-exporter

### Public Cloud Exporters

* AWS - Cloudwatch-exporter
* Azure - Promitor
* Google Cloud - Stackdriver-exporter

## Queues

* Apache Kafka
* RabbitMQ
