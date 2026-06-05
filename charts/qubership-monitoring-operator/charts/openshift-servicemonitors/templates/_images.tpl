{{- define "image" -}}
  {{- if .Values.image -}}
    {{- printf "%s" .Values.image -}}
  {{- else -}}
    {{- print "ghcr.io/netcracker/qubership-openshift-copy-certs:main" -}}
  {{- end -}}
{{- end -}}
