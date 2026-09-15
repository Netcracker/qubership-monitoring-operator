{{/* vim: set filetype=mustache: */}}

{{/*
Find a stackdriver-exporter image in various places.
Image can be found from:
* .Values.stackdriverExporter.image from values file
* or default value
*/}}
{{- define "stackdriver-exporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus-community/stackdriver_exporter */ -}}
    {{- print "docker.io/prometheuscommunity/stackdriver-exporter:v0.18.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for stackdriver-exporter.
*/}}
{{- define "stackdriver-exporter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "stackdriver-exporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
{{- end -}}
