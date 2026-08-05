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
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.imageRenderer.securityContext | default dict) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/* Return the enforced container security context for grafana-image-renderer. */}}
{{- define "grafana.imageRenderer.containerSecurityContext" -}}
{{- toYaml (dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL"))) -}}
{{- end -}}
