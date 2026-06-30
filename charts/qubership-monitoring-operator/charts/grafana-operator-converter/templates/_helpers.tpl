{{/*
Expand the name of the chart.
*/}}
{{- define "grafana-operator-converter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "grafana-operator-converter.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Allow the release namespace to be overridden
*/}}
{{- define "grafana-operator-converter.namespace" -}}
{{ .Values.namespaceOverride | default .Release.Namespace }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "grafana-operator-converter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "grafana-operator-converter.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Find a grafana-operator-converter image.
*/}}
{{- define "grafana-operator-converter.image" -}}
{{- if .Values.image -}}
{{- printf "%s" .Values.image -}}
{{- else -}}
{{- print "ghcr.io/netcracker/qubership-grafana-operator-converter:main" -}}
{{- end -}}
{{- end }}

{{- define "grafana-operator-converter.version" -}}
{{- splitList ":" (include "grafana-operator-converter.image" .) | last }}
{{- end -}}
