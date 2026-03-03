{{/*
Resource labels: name, app.kubernetes.io/name, component, part-of, managed-by, managed-by-operator.
Optional: managedByOperator (override), processedByOperator (for CRs only).
No Netcracker-specific labels (sessionId, ingress-*) — add when explicitly requested.
Usage: {{- include "common-labels.resourceLabels" (dict "ctx" . "name" $name "component" $component) | nindent 4 }}
       With override: (dict "ctx" . "name" $name "component" $component "managedByOperator" "prometheus-adapter-operator")
       For CRs: add "processedByOperator" "prometheus-adapter-operator"
       For workloads (Deployment, StatefulSet, DaemonSet), add instance, version, technology inline.
*/}}
{{- define "common-labels.resourceLabels" -}}
{{- $ctx := .ctx -}}
{{- $name := .name -}}
{{- $component := .component -}}
name: {{ $name }}
app.kubernetes.io/name: {{ $name }}
app.kubernetes.io/component: {{ $component }}
app.kubernetes.io/part-of: {{ (index $ctx.Values "global" | default dict).partOf | default "monitoring" }}
app.kubernetes.io/managed-by: {{ $ctx.Release.Service }}
app.kubernetes.io/managed-by-operator: {{ .managedByOperator | default (index $ctx.Values "global" | default dict).managedByOperator | default "monitoring-operator" }}
{{- if .processedByOperator }}
app.kubernetes.io/processed-by-operator: {{ .processedByOperator }}
{{- end }}
{{- end -}}
