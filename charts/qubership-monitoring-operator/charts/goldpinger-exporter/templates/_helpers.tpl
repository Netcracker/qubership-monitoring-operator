{{/*
Service account and image helpers for goldpinger-exporter. Labels use parent chart include "monitoring.labels" . (Values.name drives app.kubernetes.io/name).
*/}}
{{/*
Create the name of the service account to use
*/}}
{{- define "goldpinger.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default .Values.name .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Find a goldpinger-exporter image in various places.
Image can be found from:
* .Values.image from values file
* or default value
*/}}
{{- define "goldpinger.image" -}}
{{- if and .Values.image.repository .Values.image.tag -}}
{{- printf "%s:%s" .Values.image.repository .Values.image.tag -}}
{{- else -}}
{{- /* # renovate: datasource=github-releases depName=bloomberg/goldpinger */ -}}
{{- print "bloomberg/goldpinger:3.10.2" -}}
{{- end -}}
{{- end -}}
