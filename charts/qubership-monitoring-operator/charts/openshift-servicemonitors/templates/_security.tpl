{{- define "securityContext" -}}
{{- if .Values.securityContext -}}
  {{- toYaml .Values.securityContext | nindent 12 }}
{{- else if not (.Capabilities.APIVersions.Has "security.openshift.io/v1/SecurityContextConstraints") -}}
{{- toYaml (dict "runAsUser" 0 "runAsGroup" 0) | nindent 12 }}
{{- else -}}
{{- printf "{}" | nindent 12 }}
{{- end -}}
{{- end -}}