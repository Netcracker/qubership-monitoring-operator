*** Settings ***
Resource        keywords.robot

Suite Setup     Preparation Operator Session
Suite Teardown  Close Session

*** Variables ***
${vmagent_url_state}           ''
${vmsingle_url_state}          ''
${vmalert_url_state}           ''
${vmalertmanager_url_state}    ''
${vmauth_url_state}            ''

*** Test Cases ***
Check Grafana Deployment Pods Are Running
    [Tags]  full  smoke  grafana
    ${grafana_deployment_flag}=  Check Deployment State With Prerequisite  ${grafana-deployment}  ${grafana-in-cr}
    Set Suite Variable  ${grafana_deployment_flag}

Check Grafana Operator Pods Are Running
    [Tags]  full  smoke  grafana
    ${grafana_operator_flag}=  Check Deployment State With Prerequisite  ${grafana-operator}  ${grafana-in-cr}
    Set Suite Variable  ${grafana_operator_flag}

Check Monitoring Operator Pods Are Running
    [Tags]  full  smoke
    ${status}  ${monitoring_operator_flag}=  Run Keyword And Ignore Error  Check Deployment State  ${monitoring-operator}
    Run Keyword If  '${status}'=='FAIL'  Check Deployment State  ${monitoring-operator-promet}
    Set Suite Variable  ${monitoring_operator_flag}

Check Prometheus Operator Pods Are Running
    [Tags]  full  smoke-test-prometheus  smoke
    ${prometheus_operator_flag}=  Check Deployment State  ${prometheus-operator}
    Set Suite Variable  ${prometheus_operator_flag}

Check Alertmanager Pods Are Running
    [Tags]  full  smoke-test-prometheus  smoke  alertmanager
    ${alertmanager_flag}=  Check Stateful Set State With Prerequisite  ${alertmanager}-k8s  ${alertmanager-in-cr}
    Set Suite Variable  ${alertmanager_flag}

Check Prometheus Pods Are Running
    [Tags]  full  smoke-test-prometheus  smoke
    ${prometheus_flag}=  Check Stateful Set State  ${prometheus}-k8s
    Set Suite Variable  ${prometheus_flag}

Check Node Exporter Pods Are Running
    [Tags]  full  smoke
    ${node_exporter_flag}=  Check Daemon Set State With Prerequisite  ${node-exporter}  ${node-exporter-in-cr}
    Set Suite Variable  ${node_exporter_flag}

Check Kube State Metrics Pods Are Running
    [Tags]  full  smoke
    ${kube_state_metrics_flag}=  Check Deployment State With Prerequisite  ${kube-state-metrics}  ${kube-state-metrics-in-cr}
    Set Suite Variable  ${kube_state_metrics_flag}

Check Pushgateway Pods Are Running
    [Tags]  full  smoke
    ${pushgateway_flag}=  Check Deployment State With Prerequisite  ${pushgateway}  ${pushgateway-in-cr}
    Set Suite Variable  ${pushgateway_flag}

Check JSON Exporter Pods Are Running
    [Tags]  full  smoke  json-exporter
    ${json_exporter_flag}=  Check Deployment State  ${json-exporter}
    Set Suite Variable  ${json_exporter_flag}

Check Status Of Prometheus
    [Tags]  full  smoke-test-prometheus  smoke
    Check Prometheus Config Status
    Check Prometheus Flags Status
    Check Prometheus Runtime Status

Check Status Of Prometheus Api
    [Tags]  full  smoke-test-prometheus  smoke
    ${all_active_targets}=  Get All Active Targets  prometheussession
    Set Suite Variable  ${all_active_targets}
    ${prometheus_metrics}=  Get All Metrics From Api  prometheussession
    Set Suite Variable  ${prometheus_metrics}

Check Apiserver Prometheus Target Metrics
    [Tags]  full  smoke-test-prometheus  smoke
    Check Target Is UP  ${apiserver}  ${all_active_targets}
    Check Job Metrics Are Written  ${apiserver}  ${prometheus_metrics}

Check Etcd Prometheus Target Metrics
    [Tags]  full  smoke-test-prometheus  smoke  etcd
    Check Target Is UP  ${etcd}  ${all_active_targets}
    Check Job Metrics Are Written  ${etcd}  ${prometheus_metrics}

Check Kubelet Prometheus Target Metrics
    [Tags]  full  smoke-test-prometheus  smoke
    Check Target is UP  ${kubelet}  ${all_active_targets}
    Check Job Metrics Are Written  ${kubelet}  ${prometheus_metrics}

Check Non Mandatory Prometheus Target Metrics
    [Tags]  full  smoke-test-prometheus  smoke
    ${variables} =	Get Variables
    Check Kube State Metrics Target Metrics  ${kube_state_metrics_flag}  ${prometheus_metrics}
    Run Keyword If  not ${node_exporter_flag}
    ...  Fail  Please check Node Exporter pods status - now the pods don't have a Running status
    Check Node Exporter Target Metrics  ${node_exporter_flag}  ${prometheus_metrics}

Check Configurations Streamer Deployment Pods Are Running
    [Tags]  full  smoke  configurations-streamer
    ${configurations-streamer_flag}=  Check Deployment State  ${configurations-streamer}
    Set Suite Variable  ${configurations-streamer_flag}

Check Version Exporter Deployment Pods Are Running
    [Tags]  full  smoke  version-exporter
    ${version-exporter_flag}=  Check Deployment State  ${version-exporter}
    Set Suite Variable  ${version-exporter_flag}

Check Graphite Remote Adapter Deployment Pods Are Running
    [Tags]  full  smoke  graphite-remote-adapter
    ${graphite-remote-adapter_flag}=  Check Deployment State  ${graphite-remote-adapter}
    Set Suite Variable  ${graphite-remote-adapter_flag}

Check Cert Exporter Deployment Pods Are Running
    [Tags]  full  smoke  cert-exporter
    ${cert-exporter_flag}=  Check Daemon Set and Deployment State For Cert Exporter  ${cert-exporter}
    Set Suite Variable  ${cert-exporter_flag}

Check Cloudwatch Exporter Deployment Pods Are Running
    [Tags]  full  smoke  cloudwatch-exporter
    ${cloudwatch-exporter_flag}=  Check Deployment State  ${cloudwatch-exporter}
    Set Suite Variable  ${cloudwatch-exporter_flag}

Check Blackbox Exporter Pods Are Running
    [Tags]  full  smoke  blackbox-exporter
    ${blackbox-exporter_flag}=  Check Deployment Or DaemonSet State  ${blackbox-exporter}
    Set Suite Variable  ${blackbox-exporter_flag}

Check Prometheus Adapter Deployment Pods Are Running
    [Tags]  full  smoke-test-prometheus  smoke  prometheus-adapter
    ${prometheus-adapter_flag}=  Check Deployment State  ${prometheus-adapter}
    Set Suite Variable  ${prometheus-adapter_flag}

Check Prometheus Adapter Operator Deployment Pods Are Running
    [Tags]  full  smoke-test-prometheus  smoke  prometheus-adapter-operator
    ${prometheus-adapter-operator_flag}=  Check Deployment State  ${prometheus-adapter-operator}
    Set Suite Variable  ${prometheus-adapter-operator_flag}

Check Promxy Deployment Pods Are Running
    [Tags]  full  smoke  promxy
    ${promxy_flag}=  Check Deployment State  ${promxy}
    Set Suite Variable  ${promxy_flag}

Check Promitor Agent Scraper Deployment Pods Are Running
    [Tags]  full  smoke  promitor-agent-scraper
    ${promitor-agent-scraper_flag}=  Check Deployment State  ${promitor-agent-scraper}
    Set Suite Variable  ${promitor-agent-scraper_flag}

Check Network Latency Exporter Pods Are Running
    [Tags]  full  smoke  network-latency-exporter
    ${network-latency-exporter_flag}=  Check Daemon Set State  ${network-latency-exporter}
    Set Suite Variable  ${network-latency-exporter_flag}

Check Prometheus Route/Ingress Status
    [Tags]  full  smoke-test-prometheus  smoke
    ${prometheus_url_state}=  Check Route/Ingress Status  ${prometheus}  ${prometheus}  session=prometheus_ingress_session
    Set Suite Variable  ${prometheus_url_state}

Check Status Of Prometheus Web Api
    [Tags]  full  smoke-test-prometheus  smoke  ui-prometheus
    Run Keyword If  ${prometheus_url_state}
    ...  Check Service Web UI Status Via External Url  /graph  <title>Prometheus Time Series Collection and Processing Server  prometheus_ingress_session

Check AlertManager Route/Ingress Status
    [Tags]  full  smoke-test-prometheus  smoke  alertmanager
    ${alertmanager_url_state}=  Check Route/Ingress Status  ${alertmanager-in-cr}  ${alertmanager}  session=alertmanager_ingress_session
    Set Suite Variable  ${alertmanager_url_state}

Check AlertManager UI Status
    [Tags]  full  smoke-test-prometheus  smoke  ui-prometheus  alertmanager
    Run Keyword If  ${alertmanager_url_state}
    ...  Check Service Web UI Status Via External Url  /#/alerts  <title>Alertmanager</title>  alertmanager_ingress_session

Check Grafana Route/Ingress Status
    [Tags]  full  smoke  grafana
    ${grafana_url_state}=  Check Route/Ingress Status  ${grafana-in-cr}  ${grafana-in-cr}  session=grafana_ingress_session
    Set Suite Variable  ${grafana_url_state}

Check Grafana UI Status
    [Tags]  full  ui-prometheus  smoke  ui-vm  grafana
    Run Keyword If  ${grafana_url_state}
    ...  Check Service Web UI Status Via External Url  /login  <title>Grafana</title>  grafana_ingress_session

Check Pushgateway Route/Ingress Status
    [Tags]  full  smoke
    ${pushgateway_flag}=  Check Deployment State With Prerequisite  ${pushgateway}  ${pushgateway-in-cr}
    ${pushgateway_url_state}=  Check Route/Ingress Status  ${pushgateway}  ${pushgateway}  session=pushgateway_ingress_session
    Set Suite Variable  ${pushgateway_url_state}

Check Pushgateway UI Status
    [Tags]  full  ui-prometheus  smoke  ui-vm
    ${pushgateway_flag}=  Check Deployment State With Prerequisite  ${pushgateway}  ${pushgateway-in-cr}
    Run Keyword If  ${pushgateway_url_state}
    ...  Check Service Web UI Status Via External Url  /  <title>Prometheus Pushgateway</title>  pushgateway_ingress_session

Check Victoriametrics Operator Pods Are Running
    [Tags]  full  smoke-test-vm  smoke
    ${victoriametrics_flag}=  Check Deployment State With Prerequisite  ${victoriametrics-operator}  ${vm-operator-in-cr}  victoriametrics
    Set Suite Variable  ${victoriametrics_flag}

Check Vmagent Pods Are Running
    [Tags]  full  smoke-test-vm  smoke
    ${vmagent_flag}=  Check Deployment State With Prerequisite  ${vmagent}  ${vmagent-in-cr}  victoriametrics
    Set Suite Variable  ${vmagent_flag}

Check Vmsingle Pods Are Running
    [Tags]  full  smoke-test-vm  smoke
    ${vmsingle_flag}=  Check Deployment State With Prerequisite  ${vmsingle}  ${vmsingle-in-cr}  victoriametrics
    Set Suite Variable  ${vmsingle_flag}

Check Vmalert Pods Are Running
    [Tags]  full  smoke-test-vm  smoke
    ${vmalert_flag}=  Check Deployment State With Prerequisite  ${vmalert}  ${vmalert-in-cr}  victoriametrics
    Set Suite Variable  ${vmalert_flag}

Check Vmalertmanager Pods Are Running
    [Tags]  full  smoke-test-vm  smoke
    ${vmalertmanager_flag}=  Check That In CR Service Is Presented  ${vmalertmanager-in-cr}  victoriametrics
    Set Suite Variable  ${vmalertmanager_flag}

Check Vmauth Route/Ingress Status
    [Tags]  full  smoke-test-vm  smoke
    Skip If  '${vmauth}' == 'False'  VMauth is not installed. Route/Ingress doesn't exist
    ${vmauth_url_state}=  Check Route/Ingress Status  ${vmauth-in-cr}  ${vmauth-name}  victoriametrics  vmauth_ingress_session
    Set Suite Variable  ${vmauth_url_state}

Check Vmagent Route/Ingress Status
    [Tags]  full  smoke-test-vm  smoke
    Skip If  '${vmauth}' == 'True'  VMauth is installed. Route/Ingress doesn't exist
    ${vmagent_url_state}=  Run Keyword If  '${vmauth}' == 'False'
    ...  Check Route/Ingress Status  ${vmagent-in-cr}  ${vmagent}  victoriametrics  vmagent_ingress_session
    Set Suite Variable  ${vmagent_url_state}

Check Vmagent UI Status
    [Tags]  full  smoke-test-vm  smoke  ui-vm
    Skip If  '${vmauth}' == 'True'  VMauth is installed. Vmagent UI doesn't exist
    Run Keyword If  ${vmagent_url_state}
    ...  Check Service Web UI Status Via External Url  /  <h2>vmagent</h2>  vmagent_ingress_session

Check Targets UI Status
    [Tags]  full  smoke-test-vm  smoke  ui-vm
    Run Keyword If  ${vmauth_url_state}
    ...  Check Service Web UI Status Via External Url  /targets  endpoint  vmauth_ingress_session
    ...  ELSE  Check Service Web UI Status Via External Url  /targets  endpoint   vmagent_ingress_session

Check Vmsingle Route/Ingress Status
    [Tags]  full  smoke-test-vm  smoke
    Skip If  '${vmauth}' == 'True'  VMauth is installed. Route/Ingress doesn't exist
    ${vmsingle_url_state}=  Run Keyword If  '${vmauth}' == 'False'
    ...  Check Route/Ingress Status  ${vmsingle-in-cr}  ${vmsingle}  victoriametrics  vmsingle_ingress_session
    Set Suite Variable  ${vmsingle_url_state}

Check Vmsingle UI Status
    [Tags]  full  smoke-test-vm  smoke  ui-vm
    Run Keyword If  ${vmsingle_url_state}
    ...  Check Service Web UI Status Via External Url  /  <h2>Single-node VictoriaMetrics</h2>  vmsingle_ingress_session
    ...  ELSE  Check Service Web UI Status Via External Url  /  <h2>Single-node VictoriaMetrics</h2>  vmauth_ingress_session

Check Vmalert Route/Ingress Status
    [Tags]  full  smoke-test-vm  smoke
    ${vmalert_url_state}=  Run Keyword If  '${vmauth}' == 'False'
    ...  Check Route/Ingress Status  ${vmalert-in-cr}  ${vmalert}  victoriametrics  vmalert_ingress_session
    Set Suite Variable  ${vmalert_url_state}

Check Vmalert UI Status
    [Tags]  full  smoke-test-vm  smoke  ui-vm
    Run Keyword If  ${vmalert_url_state}
    ...  Check Service Web UI Status Via External Url  /  <title>vmalert - vmalert</title>  vmalert_ingress_session
    ...  ELSE  Check Service Web UI Status Via External Url  /vmalert  <title>vmalert - vmalert</title>  vmauth_ingress_session

Check Vmalertmanager Route/Ingress Status
    [Tags]  full  smoke-test-vm  smoke
    ${vmalertmanager_url_state}=  Check Route/Ingress Status  ${vmalertmanager-in-cr}  ${vmalertmanager}  victoriametrics  session=alertmanager_ingress_session
    Set Suite Variable  ${vmalertmanager_url_state}

Check Vmalertmanager UI Status
    [Tags]  full  smoke-test-vm  smoke  ui-vm
    Run Keyword If  ${vmalertmanager_url_state}
    ...  Check Service Web UI Status Via External Url  /  <title>Alertmanager</title>  alertmanager_ingress_session

Check Status Of Victoriametrics Api
    [Tags]  full  smoke-test-vm  smoke
    ${all_active_targets}=  Get All Active Targets  vmagentsession
    Set Suite Variable  ${all_active_targets}
    ${vmagent_metrics}=  Get All Metrics From Api  vmsinglessession
    Set Suite Variable  ${vmagent_metrics}

Check Apiserver Vmagent Target Metrics
    [Tags]  full  smoke-test-vm  smoke
    Check Target Is UP  ${apiserver}  ${all_active_targets}
    Check Job Metrics Are Written  ${apiserver}  ${vmagent_metrics}

Check Etcd Vmagent Target Metrics
    [Tags]  full  smoke-test-vm  smoke  etcd
    Check Target Is UP  ${etcd}  ${all_active_targets}
    Check Job Metrics Are Written  ${etcd}  ${vmagent_metrics}

Check Kubelet Vmagent Target Metrics
    [Tags]  full  smoke-test-vm  smoke
    Check Target is UP  ${kubelet}  ${all_active_targets}
    Check Job Metrics Are Written  ${kubelet}  ${vmagent_metrics}

Check Non Mandatory Victoriametrics Target Metrics
    [Tags]  full  smoke-test-vm  smoke
    ${variables} =	Get Variables
    Check Kube State Metrics Target Metrics  ${kube_state_metrics_flag}  ${vmagent_metrics}
    Run Keyword If  not ${node_exporter_flag}
    ...  Fail  Please check Node Exporter pods status - now the pods don't have a Running status
    Check Node Exporter Target Metrics  ${node_exporter_flag}  ${vmagent_metrics}
