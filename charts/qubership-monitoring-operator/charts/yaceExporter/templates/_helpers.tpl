{{/* vim: set filetype=mustache: */}}

{{/*
Find a yace-exporter image in various places.
Image can be found from:
* .Values.yaceExporter.image from values file
* or default value
*/}}
{{- define "yace-exporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus-community/yet-another-cloudwatch-exporter */ -}}
    {{- print "docker.io/prometheuscommunity/yet-another-cloudwatch-exporter:v0.67.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for yace-exporter.
*/}}
{{- define "yace-exporter.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.securityContext "id" 65534) -}}
{{- end -}}

{{/* Return the enforced container security context. */}}
{{- define "yace-exporter.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" dict) -}}
{{- end -}}

{{/*
Namespace need truncate to 26 symbols to allow specify suffixes till 35 symbols
*/}}
{{- define "monitoring.namespace" -}}
  {{- printf "%s" .Release.Namespace | trunc 26 | trimSuffix "-" -}}
{{- end -}}

{{/*
Fullname prefixed with the truncated namespace.
The namespace keeps cluster-scoped object names unique across releases.
*/}}
{{- define "yace-exporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "yace-exporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "yace-exporter.version" -}}
  {{- splitList ":" (include "yace-exporter.image" .) | last }}
{{- end -}}

{{/*
True when static AWS credentials are configured, either through chart values or an external Secret.
When only IRSA or an instance profile is used, leave aws.secret.name and both keys unset.
*/}}
{{- define "yace-exporter.awsCredentialsEnabled" -}}
{{- if or ((.Values.aws.secret).name) (and .Values.aws.aws_access_key_id .Values.aws.aws_secret_access_key) -}}true{{- end -}}
{{- end -}}

{{/*
Fail when only one of aws_access_key_id and aws_secret_access_key is set.
IRSA and aws.secret.name leave both empty, so the two keys must always be set together.
*/}}
{{- define "yace-exporter.validateAwsCredentials" -}}
  {{- if and .Values.install (not ((.Values.aws.secret).name)) -}}
    {{- $id := .Values.aws.aws_access_key_id | default "" -}}
    {{- $key := .Values.aws.aws_secret_access_key | default "" -}}
    {{- if ne (empty $id) (empty $key) -}}
      {{- fail "yace-exporter: aws.aws_access_key_id and aws.aws_secret_access_key must both be set, or both left empty when using IRSA" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}
