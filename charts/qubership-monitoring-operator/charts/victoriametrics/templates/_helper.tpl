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
    {{- /* # renovate: datasource=docker depName=rancher/kuberlr-kubectl */ -}}
    {{- print  "docker.io/rancher/kuberlr-kubectl:v8.0.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vm cleanup.
*/}}
{{- define "vm.cleanup.securityContext" -}}
{{- $legacySecurityContext := .Values.cleanup.securityContext | default dict -}}
{{- $configured := deepCopy (.Values.cleanup.hook.securityContext | default $legacySecurityContext) -}}
{{- $required := dict "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $_ := set $required "runAsNonRoot" true -}}
{{- $defaults = dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- else -}}
{{- $_ := unset $configured "runAsNonRoot" -}}{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/*
Return the enforced container security context for the VM cleanup hook.
*/}}
{{- define "vm.cleanup.containerSecurityContext" -}}
{{- $required := dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL")) -}}
{{- toYaml (mergeOverwrite (deepCopy (.Values.cleanup.hook.containerSecurityContext | default dict)) $required) -}}
{{- end -}}
