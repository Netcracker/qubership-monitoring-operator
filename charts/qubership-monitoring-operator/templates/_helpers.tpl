{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart. This is suffixed with -alertmanager, which means subtract 13 from longest 63 available
*/}}
{{- define "monitoring.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 50 | trimSuffix "-" -}}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
The components in this chart create additional resources that expand the longest created name strings.
The longest name that gets created adds and extra 37 characters, so truncation should be 63-35=26.
*/}}
{{- define "monitoring.fullname" -}}
  {{- if .Values.fullnameOverride -}}
    {{- .Values.fullnameOverride | trunc 26 | trimSuffix "-" -}}
  {{- else -}}
    {{- $name := default .Chart.Name .Values.nameOverride -}}
    {{- if contains $name .Release.Name -}}
      {{- .Release.Name | trunc 26 | trimSuffix "-" -}}
    {{- else -}}
      {{- printf "%s-%s" .Release.Name $name | trunc 26 | trimSuffix "-" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Namespace need truncate to 26 symbols to allow specify suffixes till 35 symbols
*/}}
{{- define "monitoring.namespace" -}}
  {{- printf "%s" .Release.Namespace | trunc 26 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitoring.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.monitoringOperator.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/*
Fullname suffixed with -operator
Adding 9 to 26 truncation of monitoring.fullname
*/}}
{{- define "monitoring.operator.fullname" -}}
  {{- if .Values.monitoringOperator.fullnameOverride -}}
    {{- .Values.monitoringOperator.fullnameOverride | trunc 35 | trimSuffix "-" -}}
  {{- else -}}
    {{- printf "%s-operator" (include "monitoring.fullname" .) -}}
  {{- end }}
{{- end -}}

{{/*
Fullname suffixed with -operator
Adding 9 to 26 truncation of monitoring.fullname
*/}}
{{- define "monitoring.operator.rbac.fullname" -}}
  {{- if .Values.monitoringOperator.clusterRole.name -}}
    {{- .Values.monitoringOperator.clusterRole.name | trunc 35 | trimSuffix "-" -}}
  {{- else -}}
    {{- printf "%s-%s" (include "monitoring.namespace" .) (include "monitoring.fullname" .) -}}
  {{- end }}
{{- end -}}

{{- define "monitoring.operator.version" -}}
  {{- splitList ":" (include "monitoring.operator.image" .) | last }}
{{- end -}}

{{/*
Fullname suffixed with -operator
Adding 9 to 26 truncation of monitoring.fullname
*/}}
{{- define "integrationTests.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.integrationTests.name | trunc 35 | trimSuffix "-" -}}
{{- end -}}

{{- define "integrationTests.version" -}}
  {{- splitList ":" (include "integrationTests.image" .) | last }}
{{- end -}}

{{/********************************* Kubernetes API versions ***********************************/}}

{{/* Allow KubeVersion to be overridden. */}}
{{- define "monitoring.kubeVersion" -}}
  {{- default .Capabilities.KubeVersion.Version .Values.kubeVersionOverride -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for rbac.
*/}}
{{- define "rbac.apiVersion" -}}
  {{- if semverCompare ">= 1.22-0" (include "monitoring.kubeVersion" .) -}}
    {{- print "rbac.authorization.k8s.io/v1" -}}
  {{- else -}}
    {{- print "rbac.authorization.k8s.io/v1beta1" -}}
  {{- end -}}
{{- end -}}

{{/********************************** Remote Write defaults ************************************/}}

{{/*
Return remoteWrite URLs for prometheus.
*/}}
{{- define "prometheus.remoteWrite" -}}
  {{- if not .Values.prometheus.remoteWrite -}}
    {{- if .Values.graphite_remote_adapter -}}
      {{- if .Values.graphite_remote_adapter.install -}}
      - url: "http://{{ .Values.graphite_remote_adapter.name }}:9201/write"
      {{- end -}}
    {{- else -}}
      []
    {{- end -}}
  {{- else -}}
      {{- toYaml .Values.prometheus.remoteWrite | nindent 6 }}
  {{- end -}}
{{- end -}}

{{/************************************ Ingresses *************************************/}}

{{/*
Set default value for grafana ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be used to provide Cloud DNS name.
*/}}
{{- define "grafana.ingress" -}}
{{- $ingress := .Values.grafana.ingress -}}
install: true
{{ if $ingress.host }}
host: {{ $ingress.host | quote }}
{{ else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{ else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "grafana-%s.%s" .Release.Namespace .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{- toYaml $ingress.tls | nindent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for vmSingle ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be used to provide Cloud DNS name
*/}}
{{- define "vm.single.ingress" -}}
{{- $ingress := .Values.victoriametrics.vmSingle.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "vmsingle-%s.%s" .Release.Namespace  .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for vmSelect ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.select.ingress" -}}
{{- $ingress := .Values.victoriametrics.vmCluster.vmSelectIngress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "vmselect-%s.%s" .Release.Namespace  .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for vmAgent ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.agent.ingress" -}}
{{- $ingress := .Values.victoriametrics.vmAgent.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "vmagent-%s.%s" .Release.Namespace  .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for vmAlertManager ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.alertmanager.ingress" -}}
{{- $ingress := .Values.victoriametrics.vmAlertManager.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "vmalertmanager-%s.%s" .Release.Namespace  .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for vmAlert ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.alert.ingress" -}}
{{- $ingress := .Values.victoriametrics.vmAlert.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "vmalert-%s.%s" .Release.Namespace .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for vmAuth ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.auth.ingress" -}}
{{- $ingress := .Values.victoriametrics.vmAuth.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "vmauth-%s.%s" .Release.Namespace .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for prometheus ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "prometheus.ingress" -}}
{{- $ingress := .Values.prometheus.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "prometheus-%s.%s" .Release.Namespace .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for alertManager ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "alertmanager.ingress" -}}
{{- $ingress := .Values.alertManager.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "alertmanager-%s.%s" .Release.Namespace .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end -}}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/*
Set default value for pushgateway ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "pushgateway.ingress" -}}
{{- $ingress := .Values.pushgateway.ingress -}}
install: true
{{- if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if .Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "pushgateway-%s.%s" .Release.Namespace .Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{ toYaml $ingress.tls | indent 2 }}
{{- end -}}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- end -}}

{{/********************************* Platform Monitoring Tests *********************************/}}

{{/*
Get Custom Resource plural from path in Values
*/}}
{{- define "integrationTests.plural_custom_resource" -}}
{{- printf "%v" (index (regexSplit "/" .Values.integrationTests.statusWriting.customResourcePath 5) 3) }}
{{- end -}}

{{/*
Get Custom Resource apiGroup from path in Values
*/}}
{{- define "integrationTests.apigroup_custom_resource" -}}
{{- printf "%v" (index (regexSplit "/" .Values.integrationTests.statusWriting.customResourcePath 5) 0) }}
{{- end -}}

{{/*
Build Custom Resource Path using the Helm in-built namespace parameter
*/}}
{{- define "integrationTests.customResourcePath" -}}
  {{- printf "monitoring.qubership.org/v1/%v/platformmonitorings/platformmonitoring" .Release.Namespace }}
{{- end -}}
