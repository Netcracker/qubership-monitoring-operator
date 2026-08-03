{{/* vim: set filetype=mustache: */}}

{{/*
Return securityContext for monitoring-operator.
*/}}
{{- define "monitoring.operator.securityContext" -}}
{{- $required := dict
  "runAsNonRoot" true
  "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $_ := set $defaults "runAsUser" 2000 -}}
{{- $_ := set $defaults "runAsGroup" 2000 -}}
{{- $_ := set $defaults "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.monitoringOperator.securityContext | default dict) -}}
{{- $configuredWithDefaults := mergeOverwrite $defaults $configured -}}
{{- toYaml (mergeOverwrite $configuredWithDefaults $required) -}}
{{- end -}}
{{/*
Return the container security context for the etcd-certs-to-secret job.
*/}}
{{- define "etcdCertsJob.securityContext" -}}
{{- $required := dict
  "allowPrivilegeEscalation" false
  "readOnlyRootFilesystem" true
  "capabilities" (dict "drop" (list "ALL")) -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
{{- $_ := set $required "runAsNonRoot" true -}}
{{- else -}}
{{- $_ := set $required "runAsUser" 0 -}}
{{- $_ := set $required "runAsGroup" 0 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.etcdCertsJob.securityContext | default dict) -}}
{{- toYaml (mergeOverwrite $configured $required) -}}
{{- end -}}

{{/*
Return the pod security context for etcd-certs-to-secret workloads.
The Kubernetes workload runs as root because etcd private keys are commonly readable only by root.
*/}}
{{- define "etcdCertsJob.podSecurityContext" -}}
seccompProfile:
  type: RuntimeDefault
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" }}
runAsNonRoot: true
{{- end }}
{{- end -}}
{{/*
Return securityContext for prometheus.
*/}}
{{- define "prometheus.securityContext" -}}
{{- $required := dict
  "runAsNonRoot" true
  "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $_ := set $defaults "runAsUser" 2000 -}}
{{- $_ := set $defaults "runAsGroup" 2000 -}}
{{- $_ := set $defaults "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.prometheus.securityContext | default dict) -}}
{{- $configuredWithDefaults := mergeOverwrite $defaults $configured -}}
{{- toYaml (mergeOverwrite $configuredWithDefaults $required) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for prometheus-operator.
*/}}
{{- define "prometheus.operator.securityContext" -}}
{{- $required := dict
  "runAsNonRoot" true
  "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $_ := set $defaults "runAsUser" 2000 -}}
{{- $_ := set $defaults "runAsGroup" 2000 -}}
{{- $_ := set $defaults "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.prometheus.operator.securityContext | default dict) -}}
{{- $configuredWithDefaults := mergeOverwrite $defaults $configured -}}
{{- toYaml (mergeOverwrite $configuredWithDefaults $required) | nindent 8 -}}
{{- end -}}

{{/*
Return securityContext for vmOperator.
*/}}
{{- define "vm.operator.securityContext" -}}
  {{- if .Values.victoriametrics.vmOperator.securityContext -}}
    {{- toYaml .Values.victoriametrics.vmOperator.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return containerSecurityContext for vmOperator.
*/}}
{{- define "vm.operator.containerSecurityContext" -}}
  {{- if .Values.victoriametrics.vmOperator.containerSecurityContext -}}
    {{- toYaml .Values.victoriametrics.vmOperator.containerSecurityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        runAsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vmSingle.
*/}}
{{- define "vm.single.securityContext" -}}
  {{- if .Values.victoriametrics.vmSingle.securityContext -}}
    {{- toYaml .Values.victoriametrics.vmSingle.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        runAsGroup: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vmAgent.
*/}}
{{- define "vm.agent.securityContext" -}}
  {{- if .Values.victoriametrics.vmAgent.securityContext -}}
    {{- toYaml .Values.victoriametrics.vmAgent.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vmAlertManager.
*/}}
{{- define "vm.alertmanager.securityContext" -}}
  {{- if .Values.victoriametrics.vmAlertManager.securityContext -}}
    {{- toYaml .Values.victoriametrics.vmAlertManager.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vmAlert.
*/}}
{{- define "vm.alert.securityContext" -}}
  {{- if .Values.victoriametrics.vmAlert.securityContext -}}
    {{- toYaml .Values.victoriametrics.vmAlert.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for vmAuth.
*/}}
{{- define "vm.auth.securityContext" -}}
  {{- if .Values.victoriametrics.vmAuth.securityContext -}}
    {{- toYaml .Values.victoriametrics.vmAuth.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for alertManager.
*/}}
{{- define "alertmanager.securityContext" -}}
{{- $required := dict
  "runAsNonRoot" true
  "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $_ := set $defaults "runAsUser" 2000 -}}
{{- $_ := set $defaults "runAsGroup" 2000 -}}
{{- $_ := set $defaults "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.alertManager.securityContext | default dict) -}}
{{- $configuredWithDefaults := mergeOverwrite $defaults $configured -}}
{{- toYaml (mergeOverwrite $configuredWithDefaults $required) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for grafana.
*/}}
{{- define "grafana.securityContext" -}}
  {{- if .Values.grafana.securityContext -}}
    {{- toYaml .Values.grafana.securityContext | nindent 6 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
      runAsUser: 2000
      fsGroup: 2000
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for grafana-operator.
*/}}
{{- define "grafana.operator.securityContext" -}}
  {{- if .Values.grafana.operator.securityContext -}}
    {{- toYaml .Values.grafana.operator.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for kubeStateMetrics.
*/}}
{{- define "kubeStateMetrics.securityContext" -}}
{{- $required := dict
  "runAsNonRoot" true
  "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $_ := set $defaults "runAsUser" 2000 -}}
{{- $_ := set $defaults "runAsGroup" 2000 -}}
{{- $_ := set $defaults "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := deepCopy (.Values.kubeStateMetrics.securityContext | default dict) -}}
{{- $configuredWithDefaults := mergeOverwrite $defaults $configured -}}
{{- toYaml (mergeOverwrite $configuredWithDefaults $required) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for nodeExporter.
*/}}
{{- define "nodeExporter.securityContext" -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
      {}
{{- else -}}
{{- $defaults := dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- $configured := deepCopy (.Values.nodeExporter.securityContext | default dict) -}}
{{- toYaml (mergeOverwrite $defaults $configured) | nindent 6 -}}
{{- end -}}
{{- end -}}

{{/*
Return securityContext for pushgateway.
*/}}
{{- define "pushgateway.securityContext" -}}
{{- if .Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints" -}}
      {}
{{- else -}}
{{- $defaults := dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- $configured := deepCopy (.Values.pushgateway.securityContext | default dict) -}}
{{- toYaml (mergeOverwrite $defaults $configured) | nindent 6 -}}
{{- end -}}
{{- end -}}
