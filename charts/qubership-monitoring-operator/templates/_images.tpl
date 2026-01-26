{{/* vim: set filetype=mustache: */}}

{{/************************************ Monitoring-operator ************************************/}}

{{/*
Find a monitoring-operator image in various places.
Image can be found from:
* .Values.monitoringOperator.image from values file
* or default value
*/}}
{{- define "monitoring.operator.image" -}}
  {{- if .Values.monitoringOperator.image -}}
    {{- printf "%s" .Values.monitoringOperator.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-monitoring-operator:main" -}}
  {{- end -}}
{{- end -}}

{{/****************************************** Grafana ******************************************/}}

{{/*
Find a grafana-operator image in various places.
Image can be found from:
* .Values.grafana.operator.image from values file
* or default value
*/}}
{{- define "grafana.operator.image" -}}
  {{- if .Values.grafana.operator -}}
    {{- if .Values.grafana.operator.image -}}
      {{- printf "%s" .Values.grafana.operator.image -}}
    {{- else -}}
      {{- print "quay.io/grafana-operator/grafana-operator:v4.9.0" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Find a grafana image in various places.
Image can be found from:
* .Values.grafana.image from values file
* or default value
*/}}
{{- define "grafana.image" -}}
  {{- if .Values.grafana.image -}}
    {{- printf "%s" .Values.grafana.image -}}
  {{- else -}}
    {{- print "docker.io/grafana/grafana:11.6.5" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a grafana-init-container image in various places.
Image can be found from:
* .Values.grafana.operator.initContainerImage from values file
* or default value
*/}}
{{- define "grafana.initContainer.image" -}}
  {{- if .Values.grafana.operator -}}
    {{- if .Values.grafana.operator.initContainerImage -}}
      {{- printf "%s" .Values.grafana.operator.initContainerImage -}}
    {{- else -}}
      {{- print "ghcr.io/netcracker/qubership-grafana-plugins-init:main" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/************************************* Prometheus stack **************************************/}}

{{/*
Find a prometheus-operator image in various places.
Image can be found from:
* .Values..prometheusOperator.image from values file
* or default value
*/}}
{{- define "prometheus.operator.image" -}}
  {{- if .Values.prometheus.operator.image -}}
    {{- printf "%s" .Values.prometheus.operator.image -}}
  {{- else -}}
    {{- print "quay.io/prometheus-operator/prometheus-operator:v0.80.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a prometheus-config-reloader image in various places.
Image can be found from:
* .Values.prometheus.configmapReloaderImage from values file
* or default value
*/}}
{{- define "prometheus.configReloader.image" -}}
  {{- if .Values.prometheus.configReloaderImage -}}
    {{- printf "%s" .Values.prometheus.configReloaderImage -}}
  {{- else -}}
    {{- print "quay.io/prometheus-operator/prometheus-config-reloader:v0.80.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a alertmanager image in various places.
Image can be found from:
* .Values.alertManager.image from values file
* or default value
*/}}
{{- define "alertmanager.image" -}}
  {{- if .Values.alertManager.image -}}
    {{- printf "%s" .Values.alertManager.image -}}
  {{- else -}}
    {{- print "docker.io/prom/alertmanager:v0.28.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a prometheus image in various places.
Image can be found from:
* .Values.prometheus.image from values file
* or default value
*/}}
{{- define "prometheus.image" -}}
  {{- if .Values.prometheus.image -}}
    {{- printf "%s" .Values.prometheus.image -}}
  {{- else -}}
    {{- print "docker.io/prom/prometheus:v3.2.1" -}}
  {{- end -}}
{{- end -}}

{{/*********************************** VictoriaMetrics stack ***********************************/}}

{{/*
Find a vmOperator image in various places.
Image can be found from:
* .Values.victoriametrics.vmOperator.image from values file
* or default value
*/}}
{{- define "vm.operator.image" -}}
  {{- if .Values.victoriametrics.vmOperator.image -}}
    {{- printf "%s" .Values.victoriametrics.vmOperator.image -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/operator:v0.65.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a victoriametrics operator config-reload image in various places.
Image can be found from:
* .Values.victoriametrics.vmOperator.configReloaderImage from values file
* or default value
*/}}
{{- define "vm.operator.configReloader.image" -}}
  {{- if .Values.victoriametrics.vmOperator.configReloaderImage -}}
    {{- printf "%s" .Values.victoriametrics.vmOperator.configReloaderImage -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/operator:config-reloader-v0.65.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmagent image in various places.
Image can be found from:
* .Values.victoriametrics.vmAgent.image from values file
* or default value
*/}}
{{- define "vm.agent.image" -}}
  {{- if .Values.victoriametrics.vmAgent.image -}}
    {{- printf "%s" .Values.victoriametrics.vmAgent.image -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/vmagent:v1.130.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a configmap-reload image in various places.
Image can be found from:
* .Values.victoriametrics.vmAgent.configReloadImage from values file
* or default value
*/}}
{{- define "vm.agent.configReloader.image" -}}
  {{- if .Values.victoriametrics.vmAgent.configReloadImage -}}
    {{- printf "%s" .Values.victoriametrics.vmAgent.configReloadImage -}}
  {{- else -}}
    {{- print "quay.io/prometheus-operator/prometheus-config-reloader:v0.80.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmalert image in various places.
Image can be found from:
* .Values.victoriametrics.vmAlert.image from values file
* or default value
*/}}
{{- define "vm.alert.image" -}}
  {{- if .Values.victoriametrics.vmAlert.image -}}
    {{- printf "%s" .Values.victoriametrics.vmAlert.image -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/vmalert:v1.130.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a configmap-reload image in various places.
Image can be found from:
* .Values.victoriametrics.vmAlert.configReloadImage from values file
* or default value
*/}}
{{- define "vm.alert.configReloader.image" -}}
  {{- if .Values.victoriametrics.vmAlert.configReloadImage -}}
    {{- printf "%s" .Values.victoriametrics.vmAlert.configReloadImage -}}
  {{- else -}}
    {{- print "ghcr.io/jimmidyson/configmap-reload:v0.14.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmalertmanager image in various places.
Image can be found from:
*.Values.victoriametrics.vmAlertmanager.image from values file
* or default value
*/}}
{{- define "vm.alertmanager.image" -}}
  {{- if .Values.victoriametrics.vmAlertManager.image -}}
    {{- printf "%s" .Values.victoriametrics.vmAlertManager.image -}}
  {{- else -}}
    {{- print "docker.io/prom/alertmanager:v0.28.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a configmap-reload image in various places.
Image can be found from:
* .Values.victoriametrics.vmAlertManager.configReloadImage from values file
* or default value
*/}}
{{- define "vm.alertmanager.configReloader.image" -}}
  {{- if .Values.victoriametrics.vmAlertManager.configReloadImage -}}
    {{- printf "%s" .Values.victoriametrics.vmAlertManager.configReloadImage -}}
  {{- else -}}
    {{- print "ghcr.io/jimmidyson/configmap-reload:v0.14.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmauth image in various places.
Image can be found from:
* .Values.victoriametrics.vmAuth.image from values file
* or default value
*/}}
{{- define "vm.auth.image" -}}
  {{- if .Values.victoriametrics.vmAuth.image -}}
    {{- printf "%s" .Values.victoriametrics.vmAuth.image -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/vmauth:v1.130.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmselect image in various places.
Image can be found from:
* .Values.victoriametrics.vmCluster.vmSelectImage from values file
* or default value
*/}}
{{- define "vm.select.image" -}}
  {{- if .Values.victoriametrics.vmCluster.vmSelectImage -}}
    {{- printf "%s" .Values.victoriametrics.vmCluster.vmSelectImage -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/vmselect:v1.130.0-cluster" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmstorage image in various places.
Image can be found from:
* .Values.victoriametrics.vmCluster.vmStorageImage from values file
* or default value
*/}}
{{- define "vm.storage.image" -}}
  {{- if .Values.victoriametrics.vmCluster.vmStorageImage -}}
    {{- printf "%s" .Values.victoriametrics.vmCluster.vmStorageImage -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/vmstorage:v1.130.0-cluster" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vminsert image in various places.
Image can be found from:
* .Values.victoriametrics.vmCluster.vmInsertImage from values file
* or default value
*/}}
{{- define "vm.insert.image" -}}
  {{- if .Values.victoriametrics.vmCluster.vmInsertImage -}}
    {{- printf "%s" .Values.victoriametrics.vmCluster.vmInsertImage -}}
  {{- else -}}
    {{- print "docker.io/victoriametrics/vminsert:v1.130.0-cluster" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a configmap-reload image in various places.
Image can be found from:
* .Values.victoriametrics.vmAuth.configReloadImage from values file
* or default value
*/}}
{{- define "vm.auth.configReloader.image" -}}
  {{- if .Values.victoriametrics.vmAuth.configReloadImage -}}
    {{- printf "%s" .Values.victoriametrics.vmAuth.configReloadImage -}}
  {{- else -}}
    {{- print "quay.io/prometheus-operator/prometheus-config-reloader:v0.80.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a vmsingle image in various places.
Image can be found from:
* .Values.victoriametrics.vmSingle.image from values file
* or default value
*/}}
{{- define "vm.single.image" -}}
  {{- if .Values.victoriametrics.vmSingle.image -}}
    {{- printf "%s" .Values.victoriametrics.vmSingle.image -}}
  {{- else -}}
    {{- print  "docker.io/victoriametrics/victoria-metrics:v1.130.0" -}}
  {{- end -}}
{{- end -}}

{{/***************************************** Exporters *****************************************/}}

{{/*
Find a kube-state-metrics image in various places.
Image can be found from:
* .Values.kubeStateMetrics.image from values file
* or default value
*/}}
{{- define "kubeStateMetrics.image" -}}
  {{- if .Values.kubeStateMetrics.image -}}
    {{- printf "%s" .Values.kubeStateMetrics.image -}}
  {{- else -}}
    {{- print "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.15.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a node-exporter image in various places.
Image can be found from:
* .Values.nodeExporter.image from values file
* or default value
*/}}
{{- define "nodeExporter.image" -}}
  {{- if .Values.nodeExporter.image -}}
    {{- printf "%s" .Values.nodeExporter.image -}}
  {{- else -}}
    {{- print "docker.io/prom/node-exporter:v1.9.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a pushgateway image in various places.
Image can be found from:
* .Values.pushgateway.image from values file
* or default value
*/}}
{{- define "pushgateway.image" -}}
  {{- if .Values.pushgateway.image -}}
    {{- printf "%s" .Values.pushgateway.image -}}
  {{- else -}}
    {{- print "docker.io/prom/pushgateway:v1.11.0" -}}
  {{- end -}}
{{- end -}}

{{/********************************* Platform Monitoring Tests *********************************/}}

{{/*
Find a platform_monitoring_tests image in various places.
Image can be found from:
* .Values.platform_monitoring_tests.image from values file
* or default value
*/}}
{{- define "integrationTests.image" -}}
  {{- if .Values.integrationTests.image -}}
    {{- printf "%s" .Values.integrationTests.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-monitoring-int-tests:main" -}}
  {{- end -}}
{{- end -}}

{{/********************************* OAuth2-Proxy *********************************/}}

{{/*
Find a oauth2_proxy image in various places.
Image can be found from:
* .Values.oAuthProxy.image from values file
* or default value
*/}}
{{- define "oauth2proxy.image" -}}
  {{- if .Values.oAuthProxy.image -}}
    {{- printf "%s" .Values.oAuthProxy.image -}}
  {{- else -}}
    {{- print "quay.io/oauth2-proxy/oauth2-proxy:7.8.1" -}}
  {{- end -}}
{{- end -}}
