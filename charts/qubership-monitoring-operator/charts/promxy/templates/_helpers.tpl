{{/* vim: set filetype=mustache: */}}

{{/*
Find a promxy image in various places.
Image can be found from:
* .Values.promxy.image from values file
* or default value
*/}}
{{- define "promxy.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "quay.io/jacksontj/promxy:v0.0.92" -}}
  {{- end -}}
{{- end -}}

{{/*
Find a configmap-reload image in various places.
Image can be found from:
* .Values.promxy.configmapReload.image from values file
* or default value
*/}}
{{- define "promxy-configmap-reload.image" -}}
  {{- if .Values.configmapReload.image -}}
    {{- printf "%s" .Values.configmapReload.image -}}
  {{- else -}}
    {{- print "ghcr.io/jimmidyson/configmap-reload:v0.14.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for promxy.
*/}}
{{- define "promxy.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
