{{/* vim: set filetype=mustache: */}}

{{/*
Find a graphite-remote-adapter image in various places.
Image can be found from:
* .Values.graphite_remote_adapter.image from values file
* or default value
*/}}
{{- define "graphiteRemoteAdapter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=Netcracker/qubership-graphite-remote-adapter versioning=semver */ -}}
    {{- print "ghcr.io/netcracker/qubership-graphite-remote-adapter:0.8.1" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for graphite-remote-adapter.
*/}}
{{- define "graphiteRemoteAdapter.securityContext" -}}
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
{{- define "graphiteRemoteAdapter.containerSecurityContext" -}}
{{- toYaml (dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL"))) -}}
{{- end -}}
