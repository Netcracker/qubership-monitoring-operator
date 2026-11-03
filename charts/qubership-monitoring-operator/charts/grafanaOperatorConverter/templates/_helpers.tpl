{{/*
Expand the name of the chart.
*/}}
{{- define "grafana-operator-converter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Build the normalized namespace list used by both RBAC and the controller.
global.namespaceScope wins over watchNamespaces, as it does for the Grafana operator, so a namespaced guest
release never watches namespaces that belong to the host release. The parent chart sets both global values; the
defaults apply only when the subchart is rendered on its own.
*/}}
{{- define "grafana-operator-converter.watchNamespaces" -}}
{{- $global := .Values.global | default dict -}}
{{- if dig "namespaceScope" false $global -}}
{{- .Release.Namespace -}}
{{- else if .Values.watchNamespaces -}}
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
{{- else if or .Values.namespaceScope (not (dig "privilegedRights" true $global)) -}}
{{- .Release.Namespace -}}
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
{{- $baseName := include "grafana-operator-converter.fullname" . | trunc 47 | trimSuffix "-" -}}
{{- printf "%s-leader-election" $baseName -}}
{{- end }}

{{/*
Find a converter image in various places.
Image can be found from:
* .Values.image from values file
* or default value
*/}}
{{- define "grafana-operator-converter.image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-grafana-operator-converter:main@sha256:b36881b193ea33906d35d9eccd3a61fe2f7c1788d255da5727361a1e2c0e3ade" -}}
  {{- end -}}
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

{{/*
Return container resources. User values win per key; every unset request and limit gets a default so the converter
always runs with bounded CPU, memory, and ephemeral storage.
*/}}
{{- define "grafana-operator-converter.resources" -}}
{{- $resources := deepCopy (.Values.resources | default (dict)) -}}
{{- $requests := get $resources "requests" | default (dict) -}}
{{- $limits := get $resources "limits" | default (dict) -}}
{{- $_ := set $requests "cpu" (get $requests "cpu" | default "50m") -}}
{{- $_ := set $requests "memory" (get $requests "memory" | default "64Mi") -}}
{{- $_ := set $requests "ephemeral-storage" (get $requests "ephemeral-storage" | default "16Mi") -}}
{{- $_ := set $limits "cpu" (get $limits "cpu" | default "200m") -}}
{{- $_ := set $limits "memory" (get $limits "memory" | default "256Mi") -}}
{{- $_ := set $limits "ephemeral-storage" (get $limits "ephemeral-storage" | default "256Mi") -}}
{{- $_ := set $resources "requests" $requests -}}
{{- $_ := set $resources "limits" $limits -}}
{{- toYaml $resources -}}
{{- end -}}

{{/*
Return the RBAC namespaces for a namespace-scoped converter, or an empty list for a cluster-wide one.
*/}}
{{- define "grafana-operator-converter.roleNamespaces" -}}
{{- $watchNamespaces := include "grafana-operator-converter.watchNamespaces" . -}}
{{- if $watchNamespaces -}}
{{- toJson (splitList "," $watchNamespaces) -}}
{{- else -}}
{{- toJson (list) -}}
{{- end -}}
{{- end -}}

{{/*
Return the converter rules shared by the ClusterRole and the namespaced Roles. Pass a dict with "root" and
"clusterScoped"; only a cluster-wide watch needs to read namespaces.
*/}}
{{- define "grafana-operator-converter.rules" -}}
{{- $converter := .root.Values.config -}}
{{- if and $converter.dashboard $converter.resolveGzipConfigMapRef }}
- apiGroups:
    - ""
  resources:
    - configmaps
  verbs:
    - get
    - list
    - watch
{{- end }}
{{- if $converter.dashboard }}
- apiGroups:
    - integreatly.org
  resources:
    - grafanadashboards
  verbs:
    - get
    - list
    - watch
- apiGroups:
    - grafana.integreatly.org
  resources:
    - grafanadashboards
  verbs:
    - create
    - delete
    - get
    - list
    - update
    - watch
{{- if .clusterScoped }}
- apiGroups:
    - ""
  resources:
    - namespaces
  verbs:
    - get
    - list
    - watch
{{- end }}
{{- end }}
{{- if $converter.datasource }}
- apiGroups:
    - integreatly.org
  resources:
    - grafanadatasources
  verbs:
    - list
    - watch
- apiGroups:
    - grafana.integreatly.org
  resources:
    - grafanadatasources
  verbs:
    - create
    - get
    - update
{{- end }}
{{- if $converter.folder }}
- apiGroups:
    - integreatly.org
  resources:
    - grafanafolders
  verbs:
    - list
    - watch
- apiGroups:
    - grafana.integreatly.org
  resources:
    - grafanafolders
  verbs:
    - create
    - get
    - update
{{- end }}
{{- if $converter.notification }}
- apiGroups:
    - integreatly.org
  resources:
    - grafananotificationchannels
  verbs:
    - list
    - watch
- apiGroups:
    - grafana.integreatly.org
  resources:
    - grafanacontactpoints
  verbs:
    - create
    - get
    - update
{{- end }}
{{- end -}}
