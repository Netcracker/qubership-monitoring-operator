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
  {{- if .Values.cleanup.securityContext -}}
    {{- toYaml .Values.cleanup.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return full name for mcp-victoriametrics resources.
*/}}
{{- define "vm.mcp.fullname" -}}
  {{- default "mcp-victoriametrics" .Values.mcp.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return true when mcp-victoriametrics should be rendered.
MCP requires an operator-managed backend matching mcp.vm.type. An explicit
entrypoint overrides the backend URL but does not enable standalone deployment.
*/}}
{{- define "vm.mcp.enabled" -}}
  {{- $type := default "single" .Values.mcp.vm.type -}}
  {{- $managedSingle := and (eq $type "single") .Values.vmSingle.install -}}
  {{- $managedCluster := and (eq $type "cluster") .Values.vmCluster.install -}}
  {{- if .Values.mcp.install -}}
    {{- if not (or $managedSingle $managedCluster) -}}
      {{- fail (printf "mcp.vm.type=%s requires the matching VictoriaMetrics backend to be installed" $type) -}}
    {{- end -}}
    {{- print "true" -}}
  {{- end -}}
{{- end -}}

{{/*
Return the explicit VictoriaMetrics entrypoint or the service URL for the
operator-managed backend matching mcp.vm.type.
*/}}
{{- define "vm.mcp.entrypoint" -}}
  {{- if .Values.mcp.vm.entrypoint -}}
    {{- .Values.mcp.vm.entrypoint -}}
  {{- else if eq (default "single" .Values.mcp.vm.type) "cluster" -}}
    {{- printf "http://vmselect-k8s.%s.svc:8481" .Release.Namespace -}}
  {{- else -}}
    {{- printf "http://vmsingle-k8s.%s.svc:8428" .Release.Namespace -}}
  {{- end -}}
{{- end -}}

{{/*
Find a mcp-victoriametrics image in various places.
*/}}
{{- define "vm.mcp.image" -}}
  {{- if .Values.mcp.image -}}
    {{- printf "%s" .Values.mcp.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=docker depName=ghcr.io/victoriametrics/mcp-victoriametrics */ -}}
    {{- print "ghcr.io/victoriametrics/mcp-victoriametrics:v1.20.2" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for mcp-victoriametrics.
*/}}
{{- define "vm.mcp.securityContext" -}}
  {{- if .Values.mcp.securityContext -}}
    {{- toYaml .Values.mcp.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
