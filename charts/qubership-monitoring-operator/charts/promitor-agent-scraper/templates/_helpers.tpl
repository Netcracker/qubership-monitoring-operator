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
    {{- print "ghcr.io/tomkerkhove/promitor-agent-scraper:2.13.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for promitor-agent-scraper.
*/}}
{{- define "promitor.agentScraper.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
