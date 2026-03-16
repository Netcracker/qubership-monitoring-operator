{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "ssl-exporter.name" -}}
{{- default .Chart.Name .Values.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "ssl-exporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "ssl-exporter.labels" -}}
helm.sh/chart: {{ include "ssl-exporter.chart" . }}
{{ include "ssl-exporter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/name: {{ .Values.name }}
app.kubernetes.io/component: ssl-exporter
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: monitoring
{{- end }}

{{/*
Selector labels
*/}}
{{- define "ssl-exporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "ssl-exporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "ssl-exporter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default .Values.name .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Return the appropriate image name
*/}}
{{- define "ssl-exporter.image" -}}
{{- $repo := .Values.image.repository | default "ribbybibby/ssl-exporter" }}
{{- $tag := .Values.image.tag | default "2.4.3" }}
{{- printf "%s:%s" $repo $tag }}
{{- end }}

{{/*
Generate SSL Exporter modules configuration
*/}}
{{- define "ssl-exporter.modules" -}}
{{- $modules := dict }}
{{- $defaultModules := dict }}
{{- $_ := set $defaultModules "https-selfsigned" (dict "enabled" true "timeout" "30s" "tls_config" (dict "insecure_skip_verify" true)) }}
{{- $_ := set $defaultModules "https-external" (dict "enabled" true "timeout" "30s" "tls_config" (dict "ca_file" "/etc/ssl/certs/ca-certificates.crt")) }}
{{- $_ := set $defaultModules "https-internal" (dict "enabled" true "timeout" "30s" "tls_config" (dict "ca_file" "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")) }}
{{- $_ := set $defaultModules "file" (dict "enabled" true "timeout" "30s") }}
{{- $_ := set $defaultModules "kubernetes" (dict "enabled" true "timeout" "30s") }}
{{- $_ := set $defaultModules "kubeconfig" (dict "enabled" true "timeout" "30s") }}

{{- $userModules := $.Values.modules | default dict }}

{{- range $moduleName, $defaultConfig := $defaultModules }}
{{- $userConfig := index $userModules $moduleName | default dict }}
{{- $config := merge $defaultConfig $userConfig }}
{{- if $config.enabled }}
  {{ $moduleName }}:
{{- if eq $moduleName "https-selfsigned" }}
    prober: https
{{- else if eq $moduleName "https-external" }}
    prober: https
{{- else if eq $moduleName "https-internal" }}
    prober: https
{{- else if eq $moduleName "file" }}
    prober: file
{{- else if eq $moduleName "kubernetes" }}
    prober: kubernetes
{{- else if eq $moduleName "kubeconfig" }}
    prober: kubeconfig
{{- end }}
    timeout: {{ $config.timeout }}
{{- if $config.tls_config }}
    tls_config:
{{- range $key, $value := $config.tls_config }}
{{- if kindIs "bool" $value }}
      {{ $key }}: {{ $value }}
{{- else }}
      {{ $key }}: {{ $value | quote }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
