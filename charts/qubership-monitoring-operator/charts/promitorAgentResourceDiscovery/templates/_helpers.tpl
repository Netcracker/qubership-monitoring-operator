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
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 10000 "runAsGroup" 10000 "fsGroup" 10000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.securityContext | default dict) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "promitor.agentResourceDiscovery.containerSecurityContext" -}}
{{- $configured := omit (.Values.containerSecurityContext | default dict) "enabled" -}}
{{- $required := dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL")) -}}
{{- toYaml (mergeOverwrite (deepCopy $configured) $required) -}}
{{- end -}}
