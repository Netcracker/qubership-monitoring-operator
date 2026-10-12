{{/* vim: set filetype=mustache: */}}

{{/*
Find a version-exporter image in various places.
Image can be found from:
* .Values.version-exporter.image from values file
* or default value
*/}}
{{- define "version-exporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-version-exporter versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-version-exporter:0.6.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for version-exporter.
*/}}
{{- define "version-exporter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "version-exporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" .Values.containerSecurityContext) -}}
{{- end -}}

{{/*
Namespace need truncate to 26 symbols to allow specify suffixes till 35 symbols
*/}}
{{- define "monitoring.namespace" -}}
  {{- printf "%s" .Release.Namespace | trunc 26 | trimSuffix "-" -}}
{{- end -}}

{{/*
Fullname suffixed with -operator
Adding 9 to 26 truncation of monitoring.fullname
*/}}
{{- define "version-exporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "version-exporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "version-exporter.version" -}}
  {{- splitList ":" (include "version-exporter.image" .) | last }}
{{- end -}}
