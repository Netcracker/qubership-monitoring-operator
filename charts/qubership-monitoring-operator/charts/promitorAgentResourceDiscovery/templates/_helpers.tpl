{{/* vim: set filetype=mustache: */}}

{{/*
Find a promitor-agent-resource-discovery image in various places.
Image can be found from:
* .Values.promitorAgentResourceDiscovery.image from values file
* or default value
*/}}
{{- define "promitor.agentResourceDiscovery.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=docker depName=ghcr.io/tomkerkhove/promitor-agent-resource-discovery */ -}}
    {{- print "ghcr.io/tomkerkhove/promitor-agent-resource-discovery:0.15.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for promitor-agent-resource-discovery.
*/}}
{{- define "promitor.agentResourceDiscovery.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 10000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "promitor.agentResourceDiscovery.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" (omit (.Values.containerSecurityContext | default dict) "enabled")) -}}
{{- end -}}
