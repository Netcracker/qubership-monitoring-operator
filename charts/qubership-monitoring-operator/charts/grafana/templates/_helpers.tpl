{{/* vim: set filetype=mustache: */}}

{{/*
Find a grafana-image-renderer image in various places.
Image can be found from:
* .Values.imageRenderer.image from values file
* or default value
*/}}
{{- define "grafana.imageRenderer.image" -}}
  {{- if .Values.imageRenderer.image -}}
    {{- printf "%s" .Values.imageRenderer.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=docker depName=grafana/grafana-image-renderer */ -}}
    {{- print "docker.io/grafana/grafana-image-renderer:v5.8.3" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for grafana-image-render.
*/}}
{{- define "grafana.imageRenderer.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.imageRenderer.securityContext "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context for grafana-image-renderer. */}}
{{- define "grafana.imageRenderer.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
{{- end -}}
