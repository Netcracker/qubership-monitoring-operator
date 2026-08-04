{{/* vim: set filetype=mustache: */}}

{{/*
Find a blackbox-exporter image in various places.
Image can be found from:
* .Values.blackboxExporter.image from values file
* or default value
*/}}
{{- define "blackboxExporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus/blackbox_exporter */ -}}
    {{- print "docker.io/prom/blackbox-exporter:v0.28.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for blackboxExporter.
*/}}
{{- define "blackboxExporter.securityContext" -}}
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.securityContext | default dict) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := unset $configured "runAsUser" -}}{{- $_ := unset $configured "runAsGroup" -}}{{- $_ := unset $configured "fsGroup" -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults $configured) $required) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "blackboxExporter.containerSecurityContext" -}}
{{- $required := dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL")) -}}
{{- toYaml (mergeOverwrite (deepCopy (.Values.containerSecurityContext | default dict)) $required) -}}
{{- end -}}
