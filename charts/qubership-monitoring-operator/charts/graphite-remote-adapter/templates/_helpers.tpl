{{/* vim: set filetype=mustache: */}}

{{/*
Find a graphite-remote-adapter image in various places.
Image can be found from:
* .Values.graphite_remote_adapter.image from values file
* or default value
*/}}
{{- define "graphiteRemoteAdapter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-graphite-remote-adapter:main" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for graphite-remote-adapter.
*/}}
{{- define "graphiteRemoteAdapter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
