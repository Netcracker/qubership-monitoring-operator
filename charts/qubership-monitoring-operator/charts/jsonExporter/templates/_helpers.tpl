{{/* vim: set filetype=mustache: */}}

{{/*
Find a json-exporter image in various places.
Image can be found from:
* .Values.jsonExporter.image from values file
* or default value
*/}}
{{- define "jsonExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus-community/json_exporter */ -}}
    {{- print "docker.io/prometheuscommunity/json-exporter:v0.7.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for json-exporter.
*/}}
{{- define "jsonExporter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "jsonExporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" .Values.containerSecurityContext) -}}
{{- end -}}
