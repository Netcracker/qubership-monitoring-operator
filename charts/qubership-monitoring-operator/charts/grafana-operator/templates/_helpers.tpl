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
    {{- print "docker.io/grafana/grafana-image-renderer:3.12.9" -}}
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
