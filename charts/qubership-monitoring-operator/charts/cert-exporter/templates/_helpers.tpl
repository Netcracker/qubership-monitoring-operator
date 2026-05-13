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
    {{- print "docker.io/joeelliott/cert-exporter:v2.14.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for deployment certExporter.
*/}}
{{- define "certExporter.deployment.securityContext" -}}
  {{- if .Values.deployment.securityContext -}}
    {{- toYaml .Values.deployment.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for daemonset certExporter.
*/}}
{{- define "certExporter.daemonset.securityContext" -}}
  {{- if .Values.daemonset.securityContext -}}
    {{- toYaml .Values.daemonset.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
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
