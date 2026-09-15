{{/* vim: set filetype=mustache: */}}

{{/*
Find a cert-exporter image in various places.
Image can be found from:
* .Values.certExporter.image from values file
* or default value
*/}}
{{- define "certExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=joe-elliott/cert-exporter */ -}}
    {{- print "docker.io/joeelliott/cert-exporter:v3.15.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for deployment certExporter.
*/}}
{{- define "certExporter.deployment.securityContext" -}}
{{- include "certExporter.securityContext" (dict "root" . "configured" .Values.deployment.securityContext) -}}
{{- end -}}

{{/*
Return securityContext for daemonset certExporter.
*/}}
{{- define "certExporter.daemonset.securityContext" -}}
{{- include "certExporter.securityContext" (dict "root" . "configured" .Values.daemonset.securityContext) -}}
{{- end -}}

{{/* Return the enforced pod security context. */}}
{{- define "certExporter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" .root "configured" .configured "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "certExporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
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
{{- define "certExporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "certExporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "certExporter.version" -}}
  {{- splitList ":" (include "certExporter.image" .) | last }}
{{- end -}}
