{{/* vim: set filetype=mustache: */}}

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
{{- define "vm.cleanup.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.cleanup.hook.name | trunc 35 | trimSuffix "-" -}}
{{- end -}}

{{/*
Find a vmsingle image in various places.
Image can be found from:
* .Values.cleanup.hook.image
* or default value
*/}}
{{- define "vm.cleanup.image" -}}
  {{- if .Values.cleanup.hook.image -}}
    {{- printf "%s" .Values.cleanup.hook.image -}}
  {{- else -}}
    {{- print  "docker.io/bitnami/kubectl:1.33.4" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vm cleanup.
*/}}
{{- define "vm.cleanup.securityContext" -}}
  {{- if .Values.cleanup.securityContext -}}
    {{- toYaml .Values.cleanup.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
