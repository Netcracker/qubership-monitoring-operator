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
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.root.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.configured | default dict) -}}
{{- if .root.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "certExporter.containerSecurityContext" -}}
{{- toYaml (dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL"))) -}}
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
