{{/* vim: set filetype=mustache: */}}

{{/*
Fail rendering when a configured pod security context conflicts with the enforced non-root baseline.
Pass the configured map (or nil) as the context.
*/}}
{{- define "monitoring.security.rejectPodConflicts" -}}
{{- $configured := . | default dict -}}
{{- $runAsUser := get $configured "runAsUser" -}}
{{- if and (not (kindIs "invalid" $runAsUser)) (eq (toString $runAsUser) "0") -}}
{{- fail "securityContext.runAsUser=0 conflicts with the enforced runAsNonRoot=true baseline. Use a non-root UID." -}}
{{- end -}}
{{- $runAsNonRoot := get $configured "runAsNonRoot" -}}
{{- if and (kindIs "bool" $runAsNonRoot) (not $runAsNonRoot) -}}
{{- fail "securityContext.runAsNonRoot=false conflicts with the enforced security baseline. Remove the field." -}}
{{- end -}}
{{- end -}}

{{/*
Fail rendering when a configured container security context contains a field the baseline cannot override.
Pass the configured map (or nil) as the context.
*/}}
{{- define "monitoring.security.rejectContainerConflicts" -}}
{{- $configured := . | default dict -}}
{{- if get $configured "privileged" -}}
{{- fail "containerSecurityContext.privileged=true conflicts with the enforced security baseline. Remove the field." -}}
{{- end -}}
{{- $capabilities := get $configured "capabilities" | default dict -}}
{{- if get $capabilities "add" -}}
{{- fail "containerSecurityContext.capabilities.add conflicts with the enforced drop-all baseline. Remove the added capabilities." -}}
{{- end -}}
{{- include "monitoring.security.rejectPodConflicts" $configured -}}
{{- end -}}

{{/* Enforce the pod baseline while preserving configured IDs. */}}
{{- define "monitoring.security.podContext" -}}
{{- include "monitoring.security.rejectPodConflicts" .configured -}}
{{- $required := dict "runAsNonRoot" true "seccompProfile" (dict "type" "RuntimeDefault") -}}
{{- $defaults := dict -}}
{{- if not (.root.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $id := .id | default 2000 -}}
{{- $defaults = dict "runAsUser" $id "runAsGroup" $id "fsGroup" $id -}}
{{- if hasKey . "defaults" -}}
{{- $defaults = .defaults -}}
{{- end -}}
{{- end -}}
{{- toYaml (mergeOverwrite (mergeOverwrite $defaults (deepCopy (.configured | default dict))) $required) -}}
{{- end -}}

{{/* Enforce the container baseline while preserving unrelated configured fields. */}}
{{- define "monitoring.security.containerContext" -}}
{{- include "monitoring.security.rejectContainerConflicts" .configured -}}
{{- $required := dict "allowPrivilegeEscalation" false "readOnlyRootFilesystem" true "capabilities" (dict "drop" (list "ALL")) -}}
{{- toYaml (mergeOverwrite (deepCopy (.configured | default dict)) $required) -}}
{{- end -}}

{{/* CR security contexts only support numeric identity fields. */}}
{{- define "monitoring.security.numericContext" -}}
{{- include "monitoring.security.rejectPodConflicts" .configured -}}
{{- $defaults := dict -}}
{{- if not (.root.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- $defaults = dict "runAsUser" 2000 "runAsGroup" 2000 "fsGroup" 2000 -}}
{{- end -}}
{{- $configured := pick (deepCopy (.configured | default dict)) "runAsUser" "runAsGroup" "fsGroup" -}}
{{- toYaml (mergeOverwrite $defaults $configured) -}}
{{- end -}}

{{/*
Return securityContext for monitoring-operator.
*/}}
{{- define "monitoring.operator.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.monitoringOperator.securityContext) -}}
{{- end -}}

{{/*
Return the enforced pod security context for root-chart cleanup hooks.
*/}}
{{- define "monitoring.cleanup.securityContext" -}}
{{- $values := .Values | toJson | fromJson -}}
{{- $cleanupHook := dig "victoriametrics" "cleanup" "hook" (dict) $values -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" (get $cleanupHook "securityContext")) -}}
{{- end -}}

{{/*
Return the enforced container security context for root-chart cleanup hooks.
*/}}
{{- define "monitoring.cleanup.containerSecurityContext" -}}
{{- $values := .Values | toJson | fromJson -}}
{{- $cleanupHook := dig "victoriametrics" "cleanup" "hook" (dict) $values -}}
{{- include "monitoring.security.containerContext" (dict "configured" (get $cleanupHook "containerSecurityContext")) -}}
{{- end -}}

{{/*
Return cleanup hook resources with a bounded writable layer and /tmp volume.
The hook can download a kubectl binary to /tmp when the image does not contain a compatible version.
*/}}
{{- define "monitoring.cleanup.resources" -}}
{{- $values := .Values | toJson | fromJson -}}
{{- $cleanupHook := dig "victoriametrics" "cleanup" "hook" (dict) $values -}}
{{- $resources := deepCopy (get $cleanupHook "resources" | default (dict)) -}}
{{- $limits := get $resources "limits" | default (dict) -}}
{{- $_ := set $limits "ephemeral-storage" (get $limits "ephemeral-storage" | default "100Mi") -}}
{{- $_ := set $resources "limits" $limits -}}
{{- toYaml $resources -}}
{{- end -}}

{{/*
Return the enforced pod security context for monitoring integration tests.
*/}}
{{- define "integrationTests.securityContext" -}}
{{- include "monitoring.security.podContext" (dict "root" . "configured" .Values.integrationTests.securityContext) -}}
{{- end -}}

{{/*
Return the enforced container security context for monitoring integration tests.
*/}}
{{- define "integrationTests.containerSecurityContext" -}}
{{- include "monitoring.security.containerContext" (dict "configured" .Values.integrationTests.containerSecurityContext) -}}
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
{{- include "monitoring.security.numericContext" (dict "root" . "configured" .Values.prometheus.securityContext) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for prometheus-operator.
*/}}
{{- define "prometheus.operator.securityContext" -}}
{{- include "monitoring.security.numericContext" (dict "root" . "configured" .Values.prometheus.operator.securityContext) | nindent 8 -}}
{{- end -}}

{{/*
Return securityContext for vmOperator.
*/}}
{{- define "vm.operator.securityContext" -}}
  {{- include "monitoring.security.rejectPodConflicts" .Values.victoriametrics.vmOperator.securityContext -}}
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
  {{- include "monitoring.security.rejectContainerConflicts" .Values.victoriametrics.vmOperator.containerSecurityContext -}}
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
  {{- include "monitoring.security.rejectPodConflicts" .Values.victoriametrics.vmSingle.securityContext -}}
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
  {{- include "monitoring.security.rejectPodConflicts" .Values.victoriametrics.vmAgent.securityContext -}}
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
  {{- include "monitoring.security.rejectPodConflicts" .Values.victoriametrics.vmAlertManager.securityContext -}}
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
  {{- include "monitoring.security.rejectPodConflicts" .Values.victoriametrics.vmAlert.securityContext -}}
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
  {{- include "monitoring.security.rejectPodConflicts" .Values.victoriametrics.vmAuth.securityContext -}}
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
{{- include "monitoring.security.numericContext" (dict "root" . "configured" .Values.alertManager.securityContext) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for grafana.
*/}}
{{- define "grafana.securityContext" -}}
  {{- include "monitoring.security.rejectPodConflicts" .Values.grafana.securityContext -}}
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
  {{- include "monitoring.security.rejectPodConflicts" .Values.grafana.operator.securityContext -}}
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
{{- include "monitoring.security.numericContext" (dict "root" . "configured" .Values.kubeStateMetrics.securityContext) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for nodeExporter.
*/}}
{{- define "nodeExporter.securityContext" -}}
{{- include "monitoring.security.numericContext" (dict "root" . "configured" .Values.nodeExporter.securityContext) | nindent 6 -}}
{{- end -}}

{{/*
Return securityContext for pushgateway.
*/}}
{{- define "pushgateway.securityContext" -}}
{{- include "monitoring.security.numericContext" (dict "root" . "configured" .Values.pushgateway.securityContext) | nindent 6 -}}
{{- end -}}
