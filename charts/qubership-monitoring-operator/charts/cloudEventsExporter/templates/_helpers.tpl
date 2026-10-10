{{/* vim: set filetype=mustache: */}}

{{/*
Find a cloud-events-exporter image in various places.
Image can be found from:
* .Values.cloudEventsExporter.image from values file
* or default value
*/}}
{{- define "cloudEventsExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-kube-events-reader versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-kube-events-reader:2.9.2" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for cloud-events-exporter.
*/}}
{{- define "cloudEventsExporter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 65534) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "cloudEventsExporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
{{- end -}}

{{/*
Namespace need truncate to 26 symbols to allow specify suffixes till 35 symbols
*/}}
{{- define "monitoring.namespace" -}}
  {{- printf "%s" .Release.Namespace | trunc 26 | trimSuffix "-" -}}
{{- end -}}

{{/*
Fullname suffixed with
Adding 9 to 26 truncation
*/}}
{{- define "cloudEventsExporter.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "cloudEventsExporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "cloudEventsExporter.version" -}}
  {{- splitList ":" (include "cloudEventsExporter.image" .) | last }}
{{- end -}}
