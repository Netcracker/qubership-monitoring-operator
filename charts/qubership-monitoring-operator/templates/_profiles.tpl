{{/* vim: set filetype=mustache: */}}

{{/*
Return resources for monitoring-operator by HWE profile.
*/}}
{{- define "monitoring.operator.resources" -}}
  {{- if .Values.monitoringOperator -}}
    {{- if .Values.monitoringOperator.resources -}}
      {{- toYaml .Values.monitoringOperator.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 70m
              memory: 256Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 100m
              memory: 256Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 70m
              memory: 64Mi
            limits:
              cpu: 200m
              memory: 256Mi
    {{- else -}}
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 100m
              memory: 256Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for prometheus by HWE profile.
*/}}
{{- define "prometheus.resources" -}}
  {{- if .Values.prometheus -}}
    {{- if .Values.prometheus.resources -}}
      {{- toYaml .Values.prometheus.resources | nindent 6 }}
    {{- else if eq .Values.global.profile "small" -}}
      requests:
        cpu: 1000m
        memory: 4Gi
      limits:
        cpu: 2000m
        memory: 6Gi
    {{- else if eq .Values.global.profile "medium" -}}
      requests:
        cpu: 2000m
        memory: 7Gi
      limits:
        cpu: 3500m
        memory: 12Gi
    {{- else if eq .Values.global.profile "large" -}}
      requests:
        cpu: 2500m
        memory: 15Gi
      limits:
        cpu: 4000m
        memory: 25Gi
    {{- else -}}
      requests:
        cpu: 1000m
        memory: 2Gi
      limits:
        cpu: 2000m
        memory: 3Gi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for prometheus-operator by HWE profile.
*/}}
{{- define "prometheus.operator.resources" -}}
  {{- if .Values.prometheus.operator -}}
    {{- if .Values.prometheus.operator.resources -}}
      {{- toYaml .Values.prometheus.operator.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 30m
          memory: 100Mi
        limits:
          cpu: 100m
          memory: 250Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 50m
          memory: 150Mi
        limits:
          cpu: 100m
          memory: 250Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 50m
          memory: 150Mi
        limits:
          cpu: 100m
          memory: 300Mi
    {{- else -}}
        requests:
          cpu: 50m
          memory: 50Mi
        limits:
          cpu: 100m
          memory: 250Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmOperator by HWE profile.
*/}}
{{- define "vm.operator.resources" -}}
  {{- if .Values.victoriametrics.vmOperator -}}
    {{- if .Values.victoriametrics.vmOperator.resources -}}
      {{- toYaml .Values.victoriametrics.vmOperator.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 50m
          memory: 100Mi
        limits:
          cpu: 100m
          memory: 200Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 70m
          memory: 150Mi
        limits:
          cpu: 150m
          memory: 300Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 150m
          memory: 300Mi
        limits:
          cpu: 300m
          memory: 500Mi
    {{- else -}}
        requests:
          cpu: 200m
          memory: 100Mi
        limits:
          cpu: 400m
          memory: 200Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmSingle by HWE profile.
*/}}
{{- define "vm.single.resources" -}}
  {{- if .Values.victoriametrics.vmSingle -}}
    {{- if .Values.victoriametrics.vmSingle.resources -}}
      {{- toYaml .Values.victoriametrics.vmSingle.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 300m
          memory: 1000Mi
        limits:
          cpu: 600m
          memory: 1500Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 1000m
          memory: 3000Mi
        limits:
          cpu: 1500m
          memory: 5000Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 2000m
          memory: 7000Mi
        limits:
          cpu: 3000m
          memory: 10000Mi
    {{- else -}}
        requests:
          cpu: 500m
          memory: 1000Mi
        limits:
          cpu: 1000m
          memory: 2000Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmAgent by HWE profile.
*/}}
{{- define "vm.agent.resources" -}}
  {{- if .Values.victoriametrics.vmAgent -}}
    {{- if .Values.victoriametrics.vmAgent.resources -}}
      {{- toYaml .Values.victoriametrics.vmAgent.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 500m
          memory: 512Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 500m
          memory: 512Mi
        limits:
          cpu: 1000m
          memory: 1024Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 1500m
          memory: 2048Mi
        limits:
          cpu: 2000m
          memory: 3500Mi
    {{- else -}}
        requests:
          cpu: 500m
          memory: 512Mi
        limits:
          cpu: 1000m
          memory: 1024Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmAlertManager by HWE profile.
*/}}
{{- define "vm.alertmanager.resources" -}}
  {{- if .Values.victoriametrics.vmAlertManager -}}
    {{- if .Values.victoriametrics.vmAlertManager.resources -}}
      {{- toYaml .Values.victoriametrics.vmAlertManager.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 30m
          memory: 50Mi
        limits:
          cpu: 70m
          memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 100m
          memory: 100Mi
        limits:
          cpu: 150m
          memory: 150Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 150m
          memory: 150Mi
        limits:
          cpu: 200m
          memory: 200Mi
    {{- else -}}
        requests:
          cpu: 30m
          memory: 56Mi
        limits:
          cpu: 100m
          memory: 256Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmAlert by HWE profile.
*/}}
{{- define "vm.alert.resources" -}}
  {{- if .Values.victoriametrics.vmAlert -}}
    {{- if .Values.victoriametrics.vmAlert.resources -}}
      {{- toYaml .Values.victoriametrics.vmAlert.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 50m
          memory: 150Mi
        limits:
          cpu: 100m
          memory: 200Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 100m
          memory: 250Mi
        limits:
          cpu: 150m
          memory: 400Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 250m
          memory: 400Mi
        limits:
          cpu: 400m
          memory: 700Mi
    {{- else -}}
        requests:
          cpu: 50m
          memory: 200Mi
        limits:
          cpu: 200m
          memory: 500Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmAuth by HWE profile.
*/}}
{{- define "vm.auth.resources" -}}
  {{- if .Values.victoriametrics.vmAuth -}}
    {{- if .Values.victoriametrics.vmAuth.resources -}}
      {{- toYaml .Values.victoriametrics.vmAuth.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 50m
          memory: 100Mi
        limits:
          cpu: 100m
          memory: 200Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 100m
          memory: 150Mi
        limits:
          cpu: 200m
          memory: 250Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 200m
          memory: 250Mi
        limits:
          cpu: 350m
          memory: 400Mi
    {{- else -}}
        requests:
          cpu: 50m
          memory: 200Mi
        limits:
          cpu: 200m
          memory: 500Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmSelect by HWE profile.
*/}}
{{- define "vm.select.resources" -}}
  {{- if .Values.victoriametrics.vmCluster.vmSelect -}}
    {{- if .Values.victoriametrics.vmCluster.vmSelect.resources -}}
        {{- toYaml .Values.victoriametrics.vmCluster.vmSelect.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 30m
              memory: 50Mi
            limits:
              cpu: 70m
              memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 150m
              memory: 150Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 150m
              memory: 150Mi
            limits:
              cpu: 200m
              memory: 200Mi
    {{- else -}}
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 50m
              memory: 64Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmStorage by HWE profile.
*/}}
{{- define "vm.storage.resources" -}}
  {{- if .Values.victoriametrics.vmCluster.vmStorage -}}
    {{- if .Values.victoriametrics.vmCluster.vmStorage.resources -}}
        {{- toYaml .Values.victoriametrics.vmCluster.vmStorage.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 300m
              memory: 256Mi
            limits:
              cpu: 300m
              memory: 256Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 500m
              memory: 512Mi
            limits:
              cpu: 500m
              memory: 512Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 1000m
              memory: 1024Mi
            limits:
              cpu: 1000m
              memory: 1024Mi
    {{- else -}}
            requests:
              cpu: 500m
              memory: 512Mi
            limits:
              cpu: 500m
              memory: 512Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for vmInsert by HWE profile.
*/}}
{{- define "vm.insert.resources" -}}
  {{- if .Values.victoriametrics.vmCluster.vmInsert -}}
    {{- if .Values.victoriametrics.vmCluster.vmInsert.resources -}}
        {{- toYaml .Values.victoriametrics.vmCluster.vmInsert.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 30m
              memory: 50Mi
            limits:
              cpu: 70m
              memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 150m
              memory: 150Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 150m
              memory: 150Mi
            limits:
              cpu: 200m
              memory: 200Mi
    {{- else -}}
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 50m
              memory: 64Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for alertManager by HWE profile.
*/}}
{{- define "alertmanager.resources" -}}
  {{- if .Values.alertManager -}}
    {{- if .Values.alertManager.resources -}}
      {{- toYaml .Values.alertManager.resources | nindent 6 }}
    {{- else if eq .Values.global.profile "small" -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 100m
        memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
      requests:
        cpu: 70m
        memory: 100Mi
      limits:
        cpu: 120m
        memory: 150Mi
    {{- else if eq .Values.global.profile "large" -}}
      requests:
        cpu: 150m
        memory: 200Mi
      limits:
        cpu: 200m
        memory: 300Mi
    {{- else -}}
      requests:
        cpu: 100m
        memory: 100Mi
      limits:
        cpu: 200m
        memory: 200Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for grafana by HWE profile.
*/}}
{{- define "grafana.resources" -}}
  {{- if .Values.grafana -}}
    {{- if .Values.grafana.resources -}}
      {{- toYaml .Values.grafana.resources | nindent 6 }}
    {{- else if eq .Values.global.profile "small" -}}
      requests:
        cpu: 250m
        memory: 300Mi
      limits:
        cpu: 400m
        memory: 400Mi
    {{- else if eq .Values.global.profile "medium" -}}
      requests:
        cpu: 400m
        memory: 400Mi
      limits:
        cpu: 500m
        memory: 500Mi
    {{- else if eq .Values.global.profile "large" -}}
      requests:
        cpu: 700m
        memory: 600Mi
      limits:
        cpu: 900m
        memory: 700Mi
    {{- else -}}
      requests:
        cpu: 300m
        memory: 400Mi
      limits:
        cpu: 500m
        memory: 800Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for grafana-operator by HWE profile.
*/}}
{{- define "grafana.operator.resources" -}}
  {{- if .Values.grafana.operator -}}
    {{- if .Values.grafana.operator.resources -}}
      {{- toYaml .Values.grafana.operator.resources | nindent 8 }}
    {{- else if eq .Values.global.profile "small" -}}
        requests:
          cpu: 30m
          memory: 50Mi
        limits:
          cpu: 70m
          memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
        requests:
          cpu: 50m
          memory: 150Mi
        limits:
          cpu: 100m
          memory: 250Mi
    {{- else if eq .Values.global.profile "large" -}}
        requests:
          cpu: 150m
          memory: 200Mi
        limits:
          cpu: 250m
          memory: 350Mi
    {{- else -}}
        requests:
          cpu: 50m
          memory: 50Mi
        limits:
          cpu: 100m
          memory: 100Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for grafana-image-renderer by HWE profile.
*/}}
{{- define "grafana.imageRenderer.resources" -}}
  {{- if .Values.imageRenderer -}}
    {{- if .Values.imageRenderer.resources -}}
      {{- toYaml .Values.imageRenderer.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 200m
              memory: 400Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 300m
              memory: 500Mi
            limits:
              cpu: 500m
              memory: 800Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 500m
              memory: 1000Mi
            limits:
              cpu: 800m
              memory: 2000Mi
    {{- else -}}
            requests:
              cpu: 150m
              memory: 250Mi
            limits:
              cpu: 300m
              memory: 500Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for kubeStateMetrics by HWE profile.
*/}}
{{- define "kubeStateMetrics.resources" -}}
  {{- if .Values.kubeStateMetrics -}}
    {{- if .Values.kubeStateMetrics.resources -}}
      {{- toYaml .Values.kubeStateMetrics.resources | nindent 6 }}
    {{- else if eq .Values.global.profile "small" -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 100m
        memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
      requests:
        cpu: 70m
        memory: 120Mi
      limits:
        cpu: 150m
        memory: 200Mi
    {{- else if eq .Values.global.profile "large" -}}
      requests:
        cpu: 100m
        memory: 200Mi
      limits:
        cpu: 200m
        memory: 300Mi
    {{- else -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 100m
        memory: 256Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for nodeExporter by HWE profile.
*/}}
{{- define "nodeExporter.resources" -}}
  {{- if .Values.nodeExporter -}}
    {{- if .Values.nodeExporter.resources -}}
      {{- toYaml .Values.nodeExporter.resources | nindent 6 }}
    {{- else if eq .Values.global.profile "small" -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 500m
        memory: 100Mi
    {{- else if eq .Values.global.profile "medium" -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 500m
        memory: 100Mi
    {{- else if eq .Values.global.profile "large" -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 500m
        memory: 100Mi
    {{- else -}}
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 500m
        memory: 100Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for pushgateway by HWE profile.
*/}}
{{- define "pushgateway.resources" -}}
  {{- if .Values.pushgateway -}}
    {{- if .Values.pushgateway.resources -}}
      {{- toYaml .Values.pushgateway.resources | nindent 6 }}
    {{- else if eq .Values.global.profile "small" -}}
      requests:
        cpu: 50m
        memory: 30Mi
      limits:
        cpu: 70m
        memory: 50Mi
    {{- else if eq .Values.global.profile "medium" -}}
      requests:
        cpu: 150m
        memory: 100Mi
      limits:
        cpu: 250m
        memory: 150Mi
    {{- else if eq .Values.global.profile "large" -}}
      requests:
        cpu: 250m
        memory: 150Mi
      limits:
        cpu: 400m
        memory: 250Mi
    {{- else -}}
      requests:
        cpu: 100m
        memory: 30Mi
      limits:
        cpu: 200m
        memory: 50Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for blackboxExporter by HWE profile.
*/}}
{{- define "blackboxExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 20m
              memory: 20Mi
            limits:
              cpu: 30m
              memory: 50Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 50m
              memory: 50Mi
            limits:
              cpu: 70m
              memory: 100Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 100m
              memory: 100Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else -}}
            requests:
              cpu: 50m
              memory: 50Mi
            limits:
              cpu: 100m
              memory: 300Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for certExporter daemonset by HWE profile.
*/}}
{{- define "certExporter.daemonset.resources" -}}
  {{- if .Values.daemonset -}}
    {{- if .Values.daemonset.resources -}}
      {{- toYaml .Values.daemonset.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 10m
              memory: 20Mi
            limits:
              cpu: 20m
              memory: 30Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 20m
              memory: 30Mi
            limits:
              cpu: 40m
              memory: 50Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 30m
              memory: 50Mi
            limits:
              cpu: 50m
              memory: 70Mi
    {{- else -}}
            requests:
              cpu: 10m
              memory: 25Mi
            limits:
              cpu: 20m
              memory: 50Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for certExporter deployment by HWE profile.
*/}}
{{- define "certExporter.deployment.resources" -}}
  {{- if .Values.deployment -}}
    {{- if .Values.deployment.resources -}}
      {{- toYaml .Values.deployment.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 30m
              memory: 64Mi
            limits:
              cpu: 100m
              memory: 256Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 30m
              memory: 128Mi
            limits:
              cpu: 100m
              memory: 256Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 50m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 512Mi
    {{- else -}}
            requests:
              cpu: 30m
              memory: 128Mi
            limits:
              cpu: 100m
              memory: 256Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for cloudwatchExporter by HWE profile.
*/}}
{{- define "cloudwatchExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 50m
              memory: 100Mi
            limits:
              cpu: 70m
              memory: 150Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 150Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 150m
              memory: 200Mi
            limits:
              cpu: 250m
              memory: 300Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 256Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for configurations-streamer by HWE profile.
*/}}
{{- define "configurationsStreamer.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 30m
              memory: 70Mi
            limits:
              cpu: 50m
              memory: 100Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 150Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 300Mi
            limits:
              cpu: 250m
              memory: 400Mi
  {{- else -}}
            requests:
              cpu: 50m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 200Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for graphite-remote-adapter by HWE profile.
*/}}
{{- define "graphiteRemoteAdapter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 100m
              memory: 150Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 250m
              memory: 400Mi
            limits:
              cpu: 400m
              memory: 700Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 500m
              memory: 1000Mi
            limits:
              cpu: 750m
              memory: 1500Mi
  {{- else -}}
            requests:
              cpu: 200m
              memory: 300Mi
            limits:
              cpu: 500m
              memory: 1000Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for jsonExporter by HWE profile.
*/}}
{{- define "jsonExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 50m
              memory: 100Mi
            limits:
              cpu: 70m
              memory: 150Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 150Mi
            limits:
              cpu: 150m
              memory: 200Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 250Mi
            limits:
              cpu: 300m
              memory: 350Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 100m
              memory: 128Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for networkLatencyExporter by HWE profile.
*/}}
{{- define "networkLatencyExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 70m
              memory: 100Mi
            limits:
              cpu: 150m
              memory: 200Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 150m
              memory: 200Mi
            limits:
              cpu: 250m
              memory: 300Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 250Mi
            limits:
              cpu: 300m
              memory: 350Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 256Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for prometheusAdapter by HWE profile.
*/}}
{{- define "prometheusAdapter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 4 }}
  {{- else if eq .Values.global.profile "small" -}}
    requests:
      cpu: 150m
      memory: 1000Mi
    limits:
      cpu: 250m
      memory: 2000Mi
  {{- else if eq .Values.global.profile "medium" -}}
    requests:
      cpu: 400m
      memory: 2000Mi
    limits:
      cpu: 500m
      memory: 3000Mi
  {{- else if eq .Values.global.profile "large" -}}
    requests:
      cpu: 500m
      memory: 3000Mi
    limits:
      cpu: 700m
      memory: 5000Mi
  {{- else -}}
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 200m
      memory: 384Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for prometheusAdapter operator by HWE profile.
*/}}
{{- define "prometheusAdapter.operator.resources" -}}
  {{- if .Values.operator -}}
    {{- if .Values.operator.resources -}}
      {{- toYaml .Values.operator.resources | nindent 10 }}
    {{- else if eq .Values.global.profile "small" -}}
          requests:
            cpu: 20m
            memory: 20Mi
          limits:
            cpu: 50m
            memory: 50Mi
    {{- else if eq .Values.global.profile "medium" -}}
          requests:
            cpu: 30m
            memory: 30Mi
          limits:
            cpu: 70m
            memory: 70Mi
    {{- else if eq .Values.global.profile "large" -}}
          requests:
            cpu: 50m
            memory: 30Mi
          limits:
            cpu: 100m
            memory: 100Mi
    {{- else -}}
          requests:
            cpu: 20m
            memory: 20Mi
          limits:
            cpu: 50m
            memory: 100Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for promitorAgentScraper by HWE profile.
*/}}
{{- define "promitor.agentScraper.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 70m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 150Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 150m
              memory: 200Mi
            limits:
              cpu: 200m
              memory: 250Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 250Mi
            limits:
              cpu: 400m
              memory: 500Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 256Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for promitorAgentResourceDiscovery by HWE profile.
*/}}
{{- define "promitor.agentResourceDiscovery.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 100m
              memory: 128Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 200m
              memory: 256Mi
            limits:
              cpu: 200m
              memory: 256Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 256Mi
            limits:
              cpu: 400m
              memory: 500Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 200m
              memory: 256Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for promxy by HWE profile.
*/}}
{{- define "promxy.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 50m
              memory: 100Mi
            limits:
              cpu: 100m
              memory: 150Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 300Mi
            limits:
              cpu: 250m
              memory: 350Mi
  {{- else -}}
            requests:
              cpu: 50m
              memory: 128Mi
            limits:
              cpu: 150m
              memory: 256Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for promxy configmapReload by HWE profile.
*/}}
{{- define "promxy.configmapReload.resources" -}}
  {{- if .Values.configmapReload -}}
    {{- if .Values.configmapReload.resources -}}
        {{- toYaml .Values.configmapReload.resources | nindent 12 }}
    {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 10m
              memory: 6Mi
            limits:
              cpu: 15m
              memory: 15Mi
    {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 10m
              memory: 10Mi
            limits:
              cpu: 15m
              memory: 15Mi
    {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 15m
              memory: 15Mi
            limits:
              cpu: 20m
              memory: 20Mi
    {{- else -}}
            requests:
              cpu: 5m
              memory: 3Mi
            limits:
              cpu: 10m
              memory: 20Mi
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/*
Return resources for stackdriverExporter by HWE profile.
*/}}
{{- define "stackdriverExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 50m
              memory: 70Mi
            limits:
              cpu: 100m
              memory: 150Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 150Mi
            limits:
              cpu: 150m
              memory: 200Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 250m
              memory: 300Mi
            limits:
              cpu: 350m
              memory: 400Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 100m
              memory: 128Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for versionExporter by HWE profile.
*/}}
{{- define "versionExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 150m
              memory: 250Mi
            limits:
              cpu: 200m
              memory: 300Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 300Mi
            limits:
              cpu: 300m
              memory: 400Mi
  {{- else -}}
            requests:
              cpu: 200m
              memory: 256Mi
            limits:
              cpu: 500m
              memory: 512Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for cloudEventsExporter by HWE profile.
*/}}
{{- define "cloudEventsExporter.resources" -}}
  {{- if .Values.resources -}}
      {{- toYaml .Values.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 150m
              memory: 250Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 150m
              memory: 250Mi
            limits:
              cpu: 200m
              memory: 300Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 300Mi
            limits:
              cpu: 300m
              memory: 400Mi
  {{- else -}}
            requests:
              cpu: 200m
              memory: 256Mi
            limits:
              cpu: 500m
              memory: 512Mi
  {{- end -}}
{{- end -}}

{{/*
Return resources for integrationTests by HWE profile.
*/}}
{{- define "integrationTests.resources" -}}
  {{- if .Values.integrationTests.resources -}}
      {{- toYaml .Values.integrationTests.resources | nindent 12 }}
  {{- else if eq .Values.global.profile "small" -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 300m
              memory: 400Mi
  {{- else if eq .Values.global.profile "medium" -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 300m
              memory: 400Mi
  {{- else if eq .Values.global.profile "large" -}}
            requests:
              cpu: 200m
              memory: 300Mi
            limits:
              cpu: 400m
              memory: 500Mi
  {{- else -}}
            requests:
              cpu: 100m
              memory: 200Mi
            limits:
              cpu: 300m
              memory: 400Mi
  {{- end -}}
{{- end -}}
