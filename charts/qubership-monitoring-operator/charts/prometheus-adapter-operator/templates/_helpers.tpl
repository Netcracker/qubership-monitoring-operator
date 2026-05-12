{{/* vim: set filetype=mustache: */}}

{{/*
Find a prometheus-adapter-operator image in various places.
Image can be found from:
* .Values.prometheusAdapter.operator.image from values file
* or default value
*/}}
{{- define "prometheusAdapter.operator.image" -}}
  {{- if .Values.operator.image -}}
    {{- printf "%s" .Values.operator.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-prometheus-adapter-operator versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-prometheus-adapter-operator:0.11.3" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a prometheus-adapter image in various places.
Image can be found from:
* .Values.prometheusAdapter.image from values file
* or default value
*/}}
{{- define "prometheusAdapter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-prometheus-adapter-operator versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-prometheus-adapter:0.11.3" -}}
  {{- end -}}
{{- end -}}

{{/*
Operator image tag (workload version label).
*/}}
{{- define "prometheusAdapter.operator.version" -}}
{{- splitList ":" (include "prometheusAdapter.operator.image" .) | last }}
{{- end }}

{{/*
Generate prometheusUrl for prometheus-adapter if it not defined
*/}}
{{- define "prometheusAdapter.prometheusUrl" -}}
  {{- if .Values.prometheusUrl -}}
    {{- printf "%s" (.Values.prometheusUrl) -}}
  {{- else -}}
    {{- if .Values.operator.tlsEnabled -}}
      {{- printf "https://vmsingle-k8s:8429" -}}
    {{- else -}}
      {{- printf "http://vmsingle-k8s:8429" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for prometheus-adapter.
*/}}
{{- define "prometheusAdapter.securityContext" -}}
  {{- if .Values.securityContext -}}
  {{- toYaml .Values.securityContext | nindent 4 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
    runAsUser: 2000
    fsGroup: 2000
  {{- else -}}
    {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for prometheus-adapter-operator.
*/}}
{{- define "prometheusAdapter.operator.securityContext" -}}
  {{- if .Values.operator.securityContext -}}
    {{- toYaml .Values.operator.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Namespace need truncate to 26 symbols to allow specify suffixes till 35 symbols
*/}}
{{- define "monitoring.namespace" -}}
  {{- printf "%s" .Release.Namespace | trunc 26 | trimSuffix "-" -}}
{{- end -}}

{{/*
Fullname suffixed with -operator
Adding 9 to 26 truncation of monitoring.fullname
*/}}
{{- define "prometheusAdapter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "prometheusAdapter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "prometheusAdapter.version" -}}
  {{- splitList ":" (include "prometheusAdapter.image" .) | last }}
{{- end -}}
