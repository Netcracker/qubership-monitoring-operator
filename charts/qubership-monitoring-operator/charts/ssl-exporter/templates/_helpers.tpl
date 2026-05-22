{{/* vim: set filetype=mustache: */}}

{{/*
Create the name of the service account to use
*/}}
{{- define "ssl-exporter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default .Values.name .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Find ssl-exporter image
*/}}
{{- define "ssl-exporter.image" -}}
{{- if .Values.image -}}
{{- printf "%s" .Values.image -}}
{{- else -}}
{{- print "ribbybibby/ssl-exporter:2.4.3" -}}
{{- end -}}
{{- end }}

{{/*
Generate SSL Exporter modules configuration
*/}}
{{- define "ssl-exporter.modules" -}}
{{- $modules := dict }}
{{- $defaultModules := dict }}
{{- $_ := set $defaultModules "https-selfsigned" (dict "enabled" true "timeout" "30s" "tls_config" (dict "insecure_skip_verify" true)) }}
{{- $_ := set $defaultModules "https-external" (dict "enabled" true "timeout" "30s" "tls_config" (dict "ca_file" "/etc/ssl/certs/ca-certificates.crt")) }}
{{- $_ := set $defaultModules "https-internal" (dict "enabled" true "timeout" "30s" "tls_config" (dict "ca_file" "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")) }}
{{- $_ := set $defaultModules "file" (dict "enabled" true "timeout" "30s") }}
{{- $_ := set $defaultModules "kubernetes" (dict "enabled" true "timeout" "30s") }}
{{- $_ := set $defaultModules "kubeconfig" (dict "enabled" true "timeout" "30s") }}

{{- $userModules := $.Values.modules | default dict }}

{{- range $moduleName, $defaultConfig := $defaultModules }}
{{- $userConfig := index $userModules $moduleName | default dict }}
{{- $config := merge $defaultConfig $userConfig }}
{{- if $config.enabled }}
  {{ $moduleName }}:
{{- if eq $moduleName "https-selfsigned" }}
    prober: https
{{- else if eq $moduleName "https-external" }}
    prober: https
{{- else if eq $moduleName "https-internal" }}
    prober: https
{{- else if eq $moduleName "file" }}
    prober: file
{{- else if eq $moduleName "kubernetes" }}
    prober: kubernetes
{{- else if eq $moduleName "kubeconfig" }}
    prober: kubeconfig
{{- end }}
    timeout: {{ $config.timeout }}
{{- if $config.tls_config }}
    tls_config:
{{- range $key, $value := $config.tls_config }}
{{- if kindIs "bool" $value }}
      {{ $key }}: {{ $value }}
{{- else }}
      {{ $key }}: {{ $value | quote }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
