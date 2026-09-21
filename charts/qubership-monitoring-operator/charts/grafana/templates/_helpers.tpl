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
    {{- print "docker.io/grafana/grafana-image-renderer:5.8.3" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for grafana-image-render.
*/}}
{{- define "grafana.imageRenderer.securityContext" -}}
  {{- if .Values.imageRenderer.securityContext -}}
    {{- toYaml .Values.imageRenderer.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return full name for mcp-grafana resources.
*/}}
{{- define "grafana.mcp.fullname" -}}
  {{- default "mcp-grafana" .Values.mcp.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Find a mcp-grafana image in various places.
*/}}
{{- define "grafana.mcp.image" -}}
  {{- $image := .Values.mcp.image -}}
  {{- $tag := default "0.17.2" $image.tag -}}
  {{- if $image.registry -}}
    {{- printf "%s/%s:%s" $image.registry $image.repository $tag -}}
  {{- else -}}
    {{- printf "%s:%s" $image.repository $tag -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for mcp-grafana.
*/}}
{{- define "grafana.mcp.securityContext" -}}
  {{- $platform := default "" .Values.PAAS_PLATFORM -}}
  {{- $isKubernetes := or (eq $platform "KUBERNETES") (and (eq $platform "") (not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints"))) -}}
  {{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
  {{- if $isKubernetes -}}
    {{- $_ := set $required "runAsUser" 1000 -}}
    {{- $_ := set $required "runAsGroup" 1000 -}}
  {{- end -}}
  {{- $context := mergeOverwrite (deepCopy (.Values.mcp.securityContext | default dict)) $required -}}
  {{- toYaml $context | nindent 8 }}
{{- end -}}

{{/*
Return the mandatory container securityContext for mcp-grafana, merged with
additional user settings.
*/}}
{{- define "grafana.mcp.containerSecurityContext" -}}
  {{- $platform := default "" .Values.PAAS_PLATFORM -}}
  {{- $isKubernetes := or (eq $platform "KUBERNETES") (and (eq $platform "") (not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints"))) -}}
  {{- $required := dict
        "allowPrivilegeEscalation" false
        "privileged" false
        "readOnlyRootFilesystem" true
        "runAsNonRoot" true
        "seccompProfile" (dict "type" "RuntimeDefault")
        "capabilities" (dict "add" (list) "drop" (list "ALL")) -}}
  {{- if $isKubernetes -}}
    {{- $_ := set $required "runAsUser" 1000 -}}
    {{- $_ := set $required "runAsGroup" 1000 -}}
  {{- end -}}
  {{- $context := mergeOverwrite (deepCopy (.Values.mcp.containerSecurityContext | default dict)) $required -}}
  {{- toYaml $context | nindent 12 }}
{{- end -}}
