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
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name.
*/}}
{{- define "grafana.ingress" -}}
  {{- if .Values.grafana.ingress -}}
      {{- toYaml .Values.grafana.ingress | nindent 6 -}}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
      host: "grafana-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for vmSingle ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.single.ingress" -}}
  {{- if .Values.victoriametrics.vmSingle.ingress -}}
        {{- toYaml .Values.victoriametrics.vmSingle.ingress | nindent 8 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
        host: "vmsingle-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for vmSelect ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.select.ingress" -}}
  {{- if .Values.victoriametrics.vmCluster.vmSelectIngress -}}
        {{- toYaml .Values.victoriametrics.vmCluster.vmSelectIngress | nindent 8 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
        host: "vmsingle-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for vmAgent ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.agent.ingress" -}}
  {{- if .Values.victoriametrics.vmAgent.ingress -}}
        {{- toYaml .Values.victoriametrics.vmAgent.ingress | nindent 8 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
        host: "vmagent-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for vmAlertManager ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.alertmanager.ingress" -}}
  {{- if .Values.victoriametrics.vmAlertManager.ingress -}}
        {{- toYaml .Values.victoriametrics.vmAlertManager.ingress | nindent 8 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
        host: "vmalertmanager-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for vmAlert ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.alert.ingress" -}}
  {{- if .Values.victoriametrics.vmAlert.ingress -}}
        {{- toYaml .Values.victoriametrics.vmAlert.ingress | nindent 8 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
        host: "vmalert-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for vmAuth ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "vm.auth.ingress" -}}
  {{- if .Values.victoriametrics.vmAuth.ingress -}}
        {{- toYaml .Values.victoriametrics.vmAuth.ingress | nindent 8 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
        host: "vmauth-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for prometheus ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "prometheus.ingress" -}}
  {{- if .Values.prometheus.ingress -}}
      {{- toYaml .Values.prometheus.ingress | nindent 6 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
      host: "prometheus-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for alertManager ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "alertmanager.ingress" -}}
  {{- if .Values.alertManager.ingress -}}
      {{- toYaml .Values.alertManager.ingress | nindent 6 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
      host: "alertmanager-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Set default value for pushgateway ingress host if not specify in Values.
CLOUD_PUBLIC_HOST is a parameter that can be use to provide Cloud DNS name
*/}}
{{- define "pushgateway.ingress" -}}
  {{- if .Values.pushgateway.ingress -}}
      {{- toYaml .Values.pushgateway.ingress | nindent 6 }}
  {{- else if .Values.CLOUD_PUBLIC_HOST -}}
      host: "pushgateway-{{ .Release.Namespace }}.{{ .Values.CLOUD_PUBLIC_HOST }}"
  {{- else -}}
      {}
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
  {{- printf "monitoring.qubership.org/v1alpha1/%v/platformmonitorings/platformmonitoring" .Release.Namespace }}
{{- end -}}
