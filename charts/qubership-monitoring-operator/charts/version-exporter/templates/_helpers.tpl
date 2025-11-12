{{/* vim: set filetype=mustache: */}}

{{/*
Find a version-exporter image in various places.
Image can be found from:
* .Values.version-exporter.image from values file
* or default value
*/}}
{{- define "version-exporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-version-exporter:latest" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for version-exporter.
*/}}
{{- define "version-exporter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
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
{{- define "version-exporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "version-exporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "version-exporter.version" -}}
  {{- splitList ":" (include "version-exporter.image" .) | last }}
{{- end -}}
