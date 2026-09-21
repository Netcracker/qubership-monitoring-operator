{{/* vim: set filetype=mustache: */}}

{{/*
Find a cloudwatch-exporter image in various places.
Image can be found from:
* .Values.cloudwatchExporter.image from values file
* or default value
*/}}
{{- define "cloudwatch-exporter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- /* # renovate: datasource=github-releases depName=prometheus/cloudwatch_exporter */ -}}
    {{- print "docker.io/prom/cloudwatch-exporter:v0.16.0" -}}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for cloudwatch-exporter.
*/}}
{{- define "cloudwatch-exporter.securityContext" -}}
  {{- if .Values.securityContext -}}
    {{- toYaml .Values.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 65534
        fsGroup: 65534
  {{- else -}}
        {}
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
{{- define "cloudwatch-exporter.rbac.fullname" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name -}}
{{- end -}}

{{- define "cloudwatch-exporter.instance" -}}
  {{- printf "%s-%s" (include "monitoring.namespace" .) .Values.name | nospace | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{- define "cloudwatch-exporter.version" -}}
  {{- splitList ":" (include "cloudwatch-exporter.image" .) | last }}
{{- end -}}

{{/*
True when static AWS credentials are configured — either via chart values or an external Secret.
When only IRSA / instance profile is used, leave aws.secret.name and keys unset.
*/}}
{{- define "cloudwatch-exporter.awsCredentialsEnabled" -}}
{{- if or ((.Values.aws.secret).name) (and .Values.aws.aws_access_key_id .Values.aws.aws_secret_access_key) -}}true{{- end -}}
{{- end -}}

{{/*
Name of the Secret that contains the credentials file.
*/}}
{{- define "cloudwatch-exporter.awsCredentialsSecretName" -}}
{{- if ((.Values.aws.secret).name) -}}{{ .Values.aws.secret.name }}{{- else -}}{{ .Values.name }}{{- end -}}
{{- end -}}

{{/*
Build INI content for the chart-managed Secret (aws_access_key_id / aws_secret_access_key in values).
The rendered text becomes the value of the 'credentials' key in the Secret.
*/}}
{{- define "cloudwatch-exporter.awsCredentialsIni" -}}
[default]
aws_access_key_id = {{ .Values.aws.aws_access_key_id }}
aws_secret_access_key = {{ .Values.aws.aws_secret_access_key }}
{{- end -}}
