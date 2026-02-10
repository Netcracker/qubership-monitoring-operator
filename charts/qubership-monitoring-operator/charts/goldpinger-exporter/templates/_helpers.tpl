{{/*
Expand the name of the chart.
*/}}
{{- define "goldpinger.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "goldpinger.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "goldpinger.labels" -}}
helm.sh/chart: {{ include "goldpinger.chart" . }}
{{ include "goldpinger.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/name: {{ default .Values.rbac.name (include "goldpinger.name" .) }}
app.kubernetes.io/component: goldpinger-exporter
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: monitoring
{{- end }}

{{/*
Selector labels
*/}}
{{- define "goldpinger.selectorLabels" -}}
app.kubernetes.io/name: {{ include "goldpinger.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "goldpinger.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (printf "%s-service-account" .Values.name) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Find a goldpinger-exporter image in various places.
Image can be found from:
* .Values.image from values file
* or default value
*/}}
{{- define "goldpinger.image" -}}
{{- if and .Values.image.repository .Values.image.tag -}}
{{- printf "%s:%s" .Values.image.repository .Values.image.tag -}}
{{- else -}}
{{- print "bloomberg/goldpinger:3.10.2" -}}
{{- end -}}
{{- end -}}
