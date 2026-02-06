{{- define "monitoring-pack-two.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitoring-pack-two.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.namespace" -}}
{{- default .Release.Namespace .Values.namespace -}}
{{- end -}}

{{- define "monitoring-pack-two.vmOperatorNamespace" -}}
{{- default "monitoring" .Values.vmOperatorNamespace -}}
{{- end -}}

{{- define "monitoring-pack-two.standardLabels" -}}
{{- $defaultLabels := dict "platform.monitoring.type" "exporter" -}}
{{- $customLabels := .Values.standardLabels | default dict -}}
{{- $labels := merge $defaultLabels $customLabels -}}
{{- range $k, $v := $labels }}
{{ $k }}: {{ $v | quote }}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmagent.baseName" -}}
{{- default "vmagent" .Values.vmAgent.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitoring-pack-two.vmagent.serviceAccountName" -}}
{{- if .Values.rbac.vmAgent.serviceAccountName }}
{{- .Values.rbac.vmAgent.serviceAccountName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmagent.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmagent.clusterRoleName" -}}
{{- if .Values.rbac.vmAgent.clusterRoleName }}
{{- .Values.rbac.vmAgent.clusterRoleName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-clusterrole" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmagent.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmagent.clusterRoleBindingName" -}}
{{- if .Values.rbac.vmAgent.clusterRoleBindingName }}
{{- .Values.rbac.vmAgent.clusterRoleBindingName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-clusterbinding" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmagent.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmagent.roleName" -}}
{{- if .Values.rbac.vmAgent.roleName }}
{{- .Values.rbac.vmAgent.roleName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-role" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmagent.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmagent.roleBindingName" -}}
{{- if .Values.rbac.vmAgent.roleBindingName }}
{{- .Values.rbac.vmAgent.roleBindingName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-rolebinding" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmagent.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmalert.baseName" -}}
{{- default "vmalert" .Values.vmAlert.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitoring-pack-two.vmalert.serviceAccountName" -}}
{{- if .Values.rbac.vmAlert.serviceAccountName }}
{{- .Values.rbac.vmAlert.serviceAccountName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmalert.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmalert.clusterRoleName" -}}
{{- if .Values.rbac.vmAlert.clusterRoleName }}
{{- .Values.rbac.vmAlert.clusterRoleName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-clusterrole" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmalert.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmalert.clusterRoleBindingName" -}}
{{- if .Values.rbac.vmAlert.clusterRoleBindingName }}
{{- .Values.rbac.vmAlert.clusterRoleBindingName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-clusterbinding" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmalert.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmalert.roleName" -}}
{{- if .Values.rbac.vmAlert.roleName }}
{{- .Values.rbac.vmAlert.roleName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-role" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmalert.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmalert.roleBindingName" -}}
{{- if .Values.rbac.vmAlert.roleBindingName }}
{{- .Values.rbac.vmAlert.roleBindingName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-rolebinding" (include "monitoring-pack-two.fullname" .) (include "monitoring-pack-two.vmalert.baseName" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmsingle.baseName" -}}
{{- default "k8s" .Values.vmSingle.name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "monitoring-pack-two.vmoperator.serviceAccountName" -}}
{{- if .Values.rbac.vmOperator.serviceAccountName }}
{{- .Values.rbac.vmOperator.serviceAccountName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
monitoring-victoriametrics-operator
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmoperator.clusterRoleName" -}}
{{- if .Values.rbac.vmOperator.clusterRoleName }}
{{- .Values.rbac.vmOperator.clusterRoleName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-vmoperator-clusterrole" (include "monitoring-pack-two.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.vmoperator.clusterRoleBindingName" -}}
{{- if .Values.rbac.vmOperator.clusterRoleBindingName }}
{{- .Values.rbac.vmOperator.clusterRoleBindingName | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-vmoperator-clusterbinding" (include "monitoring-pack-two.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.metadata.labels" -}}
{{- $labels := dict }}
{{- if .commonLabels }}
  {{- $labels = merge $labels .commonLabels }}
{{- end }}
{{- if .componentLabels }}
  {{- $labels = merge $labels .componentLabels }}
{{- end }}
{{- range $k, $v := $labels }}
{{ $k }}: {{ $v | quote }}
{{- end -}}
{{- end -}}

{{- define "monitoring-pack-two.metadata.annotations" -}}
{{- $annotations := dict }}
{{- if .commonAnnotations }}
  {{- $annotations = merge $annotations .commonAnnotations }}
{{- end }}
{{- if .componentAnnotations }}
  {{- $annotations = merge $annotations .componentAnnotations }}
{{- end }}
{{- if $annotations }}
{{- range $k, $v := $annotations }}
{{ $k }}: {{ $v | quote }}
{{- end }}
{{- end -}}
{{- end -}}
