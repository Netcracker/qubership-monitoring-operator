{{/* vim: set filetype=mustache: */}}

{{/*
Find a network-latency-exporter image in various places.
Image can be found from:
* .Values.networkLatencyExporter.image from values file
* or default value
*/}}
{{- define "networkLatencyExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-network-latency-exporter versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-network-latency-exporter:2.10.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for network-latency-exporter.
*/}}
{{- define "networkLatencyExporter.securityContext" -}}
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 2001 "runAsGroup" 2001 "fsGroup" 2001 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.securityContext | default dict) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/*
Return the enforced container security context.
Privileged mode remains an explicit compatibility override.
*/}}
{{- define "networkLatencyExporter.containerSecurityContext" -}}
{{- if .Values.rbac.privileged -}}
allowPrivilegeEscalation: true
privileged: true
readOnlyRootFilesystem: true
{{- else -}}
allowPrivilegeEscalation: false
readOnlyRootFilesystem: true
capabilities:
  drop:
    - ALL
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
{{- define "networkLatencyExporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "networkLatencyExporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "networkLatencyExporter.version" -}}
  {{- splitList ":" (include "networkLatencyExporter.image" .) | last }}
{{- end -}}
