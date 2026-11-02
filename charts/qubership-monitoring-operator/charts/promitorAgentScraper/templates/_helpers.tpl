{{/* vim: set filetype=mustache: */}}

{{/*
Find a promitor-agent-scraper image in various places.
Image can be found from:
* .Values.promitorAgentScraper.image from values file
* or default value
*/}}
{{- define "promitor.agentScraper.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=docker depName=ghcr.io/tomkerkhove/promitor-agent-scraper */ -}}
    {{- print "ghcr.io/tomkerkhove/promitor-agent-scraper:2.15.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for promitor-agent-scraper.
*/}}
{{- define "promitor.agentScraper.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "promitor.agentScraper.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
{{- end -}}
