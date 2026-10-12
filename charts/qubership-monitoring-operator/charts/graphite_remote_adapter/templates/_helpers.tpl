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
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-graphite-remote-adapter versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-graphite-remote-adapter:0.8.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for graphite-remote-adapter.
*/}}
{{- define "graphiteRemoteAdapter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "graphiteRemoteAdapter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
{{- end -}}
