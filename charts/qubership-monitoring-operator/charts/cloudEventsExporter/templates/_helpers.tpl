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
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 65534 "runAsGroup" 65534 "fsGroup" 65534 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.securityContext | default dict) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "cloudEventsExporter.containerSecurityContext" -}}
{{- toYaml (dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL"))) -}}
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
