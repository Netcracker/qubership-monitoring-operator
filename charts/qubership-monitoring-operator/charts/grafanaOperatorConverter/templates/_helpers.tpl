{{/*
Expand the name of the chart.
*/}}
{{- define "grafana-operator-converter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Build the normalized namespace list used by both RBAC and the controller.
*/}}
{{- define "grafana-operator-converter.watchNamespaces" -}}
{{- if .Values.watchNamespaces -}}
{{- $namespaces := list -}}
{{- range splitList "," .Values.watchNamespaces -}}
{{- $namespace := trim . -}}
{{- if $namespace -}}
{{- $namespaces = append $namespaces $namespace -}}
{{- end -}}
{{- end -}}
{{- if eq (len $namespaces) 0 -}}
{{- fail "watchNamespaces must contain at least one non-empty namespace" -}}
{{- end -}}
{{- join "," (sortAlpha (uniq $namespaces)) -}}
{{- else if or .Values.namespaceScope (not .Values.global.privilegedRights) -}}
{{- include "grafana-operator-converter.namespace" . -}}
{{- end -}}
{{- end -}}

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
Create names for resources that append a suffix to the full name.
*/}}
{{- define "grafana-operator-converter.metricsServiceName" -}}
{{- printf "%s-metrics-service" (include "grafana-operator-converter.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "grafana-operator-converter.leaderElectionName" -}}
{{- printf "%s-leader-election" (include "grafana-operator-converter.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{/*
Allow the release namespace to be overridden
*/}}
{{- define "grafana-operator-converter.namespace" -}}
{{ .Values.namespaceOverride | default .Release.Namespace }}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "grafana-operator-converter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "grafana-operator-converter.labels" -}}
helm.sh/chart: {{ include "grafana-operator-converter.chart" . }}
{{ include "grafana-operator-converter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "grafana-operator-converter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "grafana-operator-converter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

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
