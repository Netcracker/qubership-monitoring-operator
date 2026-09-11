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
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.securityContext | default dict) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "promitor.agentScraper.containerSecurityContext" -}}
{{- toYaml (dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL"))) -}}
{{- end -}}
