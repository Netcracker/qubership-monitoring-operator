{{/* vim: set filetype=mustache: */}}

{{/*
Base resource labels: expects either chart context (.) or dict with ctx and optional overrides.
When name or component are omitted, they are derived from ctx.Values (name, component|default name).
Usage:
  {{- include "monitoring.labels" . | nindent 4 }}
  Or with overrides: (dict "ctx" . "name" "x" "component" "y")
*/}}
{{- define "monitoring.labels" -}}
{{- $ctx := index . "ctx" | default . -}}
{{- $vals := $ctx.Values -}}
{{- $name := .name | default $vals.name -}}
{{- $component := .component | default $ctx.Chart.Name -}}
name: {{ $name }}
app.kubernetes.io/name: {{ $name }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/part-of: monitoring
app.kubernetes.io/managed-by: {{ $ctx.Release.Service }}
{{- end -}}

{{/*
Extra labels helper: renders arbitrary labels map when provided.
Usage:
  {{- include "monitoring.extraLabels" (dict "ctx" . "extraLabels" .Values.labels) | nindent 4 }}
*/}}
{{- define "monitoring.extraLabels" -}}
{{- $ctx := index . "ctx" | default . -}}
{{- $vals := $ctx.Values -}}
{{- $extra := .extraLabels | default ($vals.labels | default dict) -}}
{{ with $extra }}

{{ toYaml . }}
{{- end }}
{{- end -}}

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
Render ingress settings for PlatformMonitoring CR.
*/}}
{{- define "monitoring.ingress" -}}
{{- $ingress := .ingress -}}
{{- $ctx := .ctx -}}
{{- $hostPrefix := .hostPrefix -}}
install: true
{{- if $ingress.rules }}
rules:
{{ toYaml $ingress.rules }}
{{- else if $ingress.host }}
host: {{ $ingress.host | quote }}
{{- else if $ctx.Values.CLOUD_PUBLIC_HOST }}
host: {{ printf "%s-%s.%s" $hostPrefix $ctx.Release.Namespace $ctx.Values.CLOUD_PUBLIC_HOST | quote }}
{{- end -}}
{{- if $ingress.tls }}
tls:
{{- toYaml $ingress.tls | nindent 2 }}
{{- end }}
{{- if $ingress.tlsSecretName }}
tlsSecretName: {{ $ingress.tlsSecretName | quote }}
{{- end -}}
{{- if $ingress.ingressClassName }}
ingressClassName: {{ $ingress.ingressClassName | quote }}
{{- end }}
{{- if $ingress.annotations }}
annotations:
{{ toYaml $ingress.annotations | indent 2 }}
{{- end }}
{{- if $ingress.labels }}
labels:
{{ toYaml $ingress.labels | indent 2 }}
{{- end }}
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
  {{- printf "monitoring.netcracker.com/v1/%v/platformmonitorings/platformmonitoring" .Release.Namespace }}
{{- end -}}
