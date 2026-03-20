{{/* vim: set filetype=mustache: */}}

{{/*
Find a blackbox-exporter image in various places.
Image can be found from:
* .Values.blackboxExporter.image from values file
* or default value
*/}}
{{- define "blackboxExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "docker.io/prom/blackbox-exporter:v0.27.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for blackboxExporter.
*/}}
{{- define "blackboxExporter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
    {{- else -}}
        {}
    {{- end -}}
{{- end -}}
