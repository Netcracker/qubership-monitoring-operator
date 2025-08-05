{{/* vim: set filetype=mustache: */}}

{{/*
Return securityContext for monitoring-operator.
*/}}
{{- define "monitoring.operator.securityContext" -}}
  {{- if .Values.monitoringOperator.securityContext -}}
    {{- toYaml .Values.monitoringOperator.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
{{- end -}}
{{/*
Return securityContext for etcd-certs-to-secret job.
*/}}
{{- define "etcdCertsJob.securityContext" -}}
{{- if .Values.etcdCertsJob.securityContext -}}
  {{- toYaml .Values.etcdCertsJob.securityContext | nindent 8 }}
{{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 0
        runAsGroup: 0
{{- else -}}
            {}
{{- end -}}
{{- end -}}
{{/*
Return securityContext for etcd-certs-to-secret job.
*/}}
{{- define "etcdCertsCronJob.securityContext" -}}
{{- if .Values.etcdCertsJob.securityContext -}}
  {{- toYaml .Values.etcdCertsJob.securityContext | nindent 12 }}
{{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
            runAsUser: 0
            runAsGroup: 0
{{- else -}}
            {}
{{- end -}}
{{- end -}}
{{/*
Return securityContext for prometheus.
*/}}
{{- define "prometheus.securityContext" -}}
  {{- if .Values.prometheus.securityContext -}}
    {{- toYaml .Values.prometheus.securityContext | nindent 6 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
      runAsUser: 2000
      fsGroup: 2000
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for prometheus-operator.
*/}}
{{- define "prometheus.operator.securityContext" -}}
  {{- if .Values.prometheus.operator.securityContext -}}
    {{- toYaml .Values.prometheus.operator.securityContext | nindent 8 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
        runAsUser: 2000
        fsGroup: 2000
  {{- else -}}
        {}
  {{- end -}}
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
  {{- if .Values.alertManager.securityContext -}}
    {{- toYaml .Values.alertManager.securityContext | nindent 6 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
      runAsUser: 2000
      fsGroup: 2000
  {{- else -}}
      {}
  {{- end -}}
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
  {{- if .Values.kubeStateMetrics.securityContext -}}
    {{- toYaml .Values.kubeStateMetrics.securityContext | nindent 6 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
      runAsUser: 2000
      fsGroup: 2000
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for nodeExporter.
*/}}
{{- define "nodeExporter.securityContext" -}}
  {{- if .Values.nodeExporter.securityContext -}}
    {{- toYaml .Values.nodeExporter.securityContext | nindent 6 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
      runAsUser: 2000
      fsGroup: 2000
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}

{{/*
Return securityContext for pushgateway.
*/}}
{{- define "pushgateway.securityContext" -}}
  {{- if .Values.pushgateway.securityContext -}}
    {{- toYaml .Values.pushgateway.securityContext | nindent 6 }}
  {{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
      runAsUser: 2000
      fsGroup: 2000
  {{- else -}}
      {}
  {{- end -}}
{{- end -}}
