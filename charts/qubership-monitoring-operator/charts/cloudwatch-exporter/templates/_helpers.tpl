{{/* vim: set filetype=mustache: */}}

{{/*
Find a cloudwatch-exporter image in various places.
Image can be found from:
* .Values.cloudwatchExporter.image from values file
* or default value
*/}}
{{- define "cloudwatch-exporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus/cloudwatch_exporter */ -}}
    {{- print "docker.io/prom/cloudwatch-exporter:v0.16.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for cloudwatch-exporter.
*/}}
{{- define "cloudwatch-exporter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 65534
        fsGroup: 65534
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
{{- define "cloudwatch-exporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "cloudwatch-exporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "cloudwatch-exporter.version" -}}
  {{- splitList ":" (include "cloudwatch-exporter.image" .) | last }}
{{- end -}}
