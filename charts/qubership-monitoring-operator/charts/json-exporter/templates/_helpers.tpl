{{/* vim: set filetype=mustache: */}}

{{/*
Find a json-exporter image in various places.
Image can be found from:
* .Values.jsonExporter.image from values file
* or default value
*/}}
{{- define "jsonExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus-community/json_exporter */ -}}
    {{- print "docker.io/prometheuscommunity/json-exporter:v0.7.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for json-exporter.
*/}}
{{- define "jsonExporter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
