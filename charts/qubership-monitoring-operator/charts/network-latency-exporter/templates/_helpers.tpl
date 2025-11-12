{{/* vim: set filetype=mustache: */}}

{{/*
Find a network-latency-exporter image in various places.
Image can be found from:
* .Values.networkLatencyExporter.image from values file
* or default value
*/}}
{{- define "networkLatencyExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-network-latency-exporter:latest" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for network-latency-exporter.
*/}}
{{- define "networkLatencyExporter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 0
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
{{- define "networkLatencyExporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "networkLatencyExporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "networkLatencyExporter.version" -}}
  {{- splitList ":" (include "networkLatencyExporter.image" .) | last }}
{{- end -}}
