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
    {{- print "ghcr.io/tomkerkhove/promitor-agent-resource-discovery:0.13.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for promitor-agent-resource-discovery.
*/}}
{{- define "promitor.agentResourceDiscovery.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 10000
        runAsGroup: 10000
        runAsNonRoot: true
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
