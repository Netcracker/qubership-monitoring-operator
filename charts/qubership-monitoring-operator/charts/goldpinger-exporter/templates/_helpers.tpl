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
