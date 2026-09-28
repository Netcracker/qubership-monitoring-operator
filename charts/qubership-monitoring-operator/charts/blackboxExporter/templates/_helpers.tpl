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
{{- $configured := deepCopy (.Values.securityContext | default dict) -}}
{{- $sysctls := get $configured "sysctls" | default list -}}
{{- $pingGroupRangeConfigured := false -}}
{{- range $sysctl := $sysctls -}}
  {{- if eq (get $sysctl "name") "net.ipv4.ping_group_range" -}}
    {{- $pingGroupRangeConfigured = true -}}
  {{- end -}}
{{- end -}}
{{- if not $pingGroupRangeConfigured -}}
  {{- $sysctls = append $sysctls (dict "name" "net.ipv4.ping_group_range" "value" "0 2147483647") -}}
{{- end -}}
{{- $_ := set $configured "sysctls" $sysctls -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" $configured "id" 2000) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "blackboxExporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" .Values.containerSecurityContext) -}}
{{- end -}}
