# Useful links
# API kuber with python https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/AppsV1Api.md

*** Settings ***
Resource                        ../keywords.robot

*** Variables ***
${namespace}                    %{NAMESPACE}
${monitoring-operator}          monitoring-operator
${monitoring-operator-promet}   monitoring-operator-promet
${prometheus-operator}          prometheus-operator
${victoriametrics-operator}     victoriametrics-operator
${vm-operator-in-cr}            vmOperator
${grafana-operator}             grafana-operator
${grafana-deployment}           grafana-deployment
${grafana-in-cr}                grafana

${kube-state-metrics-in-cr}     kubeStateMetrics
${kube-state-metrics}           kube-state-metrics
${node-exporter-in-cr}          nodeExporter
${node-exporter}                node-exporter

${alertmanager}                 alertmanager
${alertmanager-in-cr}           alertManager
${prometheus}                   prometheus

${victoriametrics}              victoriametrics
${vmsingle}                     vmsingle-k8s
${vmagent}                      vmagent-k8s
${vmalert}                      vmalert-k8s
${vmalertmanager}               vmalertmanager-k8s
${vmauth-name}                  vmauth-k8s
${vmsingle-in-cr}               vmSingle
${vmagent-in-cr}                vmAgent
${vmalert-in-cr}                vmAlert
${vmalertmanager-in-cr}         vmAlertManager
${vmauth-in-cr}                 vmAuth

${pushgateway}                  pushgateway
${pushgateway-in-cr}            pushgateway

${apiserver}                    apiserver
${etcd}                         etcd
${kube-controller-manager}      kube-controller-manager
${kube-scheduler}               kube-scheduler
${kubelet}                      kubelet

${configurations-streamer}      configurations-streamer
${version-exporter}             version-exporter
${graphite-remote-adapter}      graphite-remote-adapter
${cert-exporter}                cert-exporter
${cloudwatch-exporter}          cloudwatch-exporter
${blackbox-exporter}            blackbox-exporter
${prometheus-adapter}           prometheus-adapter
${prometheus-adapter-operator}  prometheus-adapter-operator
${promxy}                       promxy
${promitor-agent-scraper}       promitor-agent-scraper
${network-latency-exporter}     network-latency-exporter
${json-exporter}                json-exporter

${FILES_PATH}                   integration-tests/source_files
${RETRY_TIME}                   5min
${RETRY_INTERVAL}               3s
${TAGS}                         %{TAGS}

*** Keywords ***
Close Session
    Delete All Sessions

Preparation Session For External Service
    [Arguments]  ${external_url}    ${session}=external
    ${auth}=    Run Keyword If  '${session}'=='vmauth_ingress_session'  Get Creadentials From Secret
    Create Session  ${session}  ${external_url}  auth=${auth}

Check Pod's List Is Not Empty
    [Arguments]  ${list_of_pods}
    ${list_len}=  Get List Length  ${list_of_pods}
    Run Keyword If  ${list_len} == 0  Fail  Error! Found zero pods!
    Should Not Be Empty  ${list_of_pods}

Check Pod's List Is Equals
    [Arguments]  ${list_of_pods}  ${desired_count_of_pods}
    ${list_len}=  Get List Length  ${list_of_pods}
    Check pod's list is not empty  ${list_of_pods}
    Run Keyword If  ${list_len} != ${desired_count_of_pods}  Fail
    ...  Error! Found ${list_len}, but expected ${desired_count_of_pods} pods!

Check Status Of Pods
    [Arguments]  ${list_pods}
    FOR  ${pod}  IN  @{list_pods}
       ${state}=  Run Keyword And Return Status  Should Be Equal As Strings  ${pod.status.phase}  Running
       Should Be True  ${state}
       ...  Error! Following pod ${pod.metadata.name} has Failed status! Please, recheck pod status
    END
    RETURN  ${state}

Determine Deployment Type
    [Arguments]  ${name}
    ${deployment_exists}=  Run Keyword And Return Status  Get Deployment Entity  ${name}  ${namespace}
    ${daemonset_exists}=   Run Keyword And Return Status  Get Daemon Set  ${name}  ${namespace}
    ${deployment_type}=    Set Variable  none
    ${deployment_type}=    Run Keyword If    ${deployment_exists} and not ${daemonset_exists}  
    ...    Set Variable    deployment  
    ...    ELSE    Set Variable    daemonset
    RETURN  ${deployment_type}
    
Check Deployment Or DaemonSet State
    [Arguments]  ${name}
    ${deployment_type}=  Determine Deployment Type  ${name}
    Run Keyword If  "${deployment_type}" == "deployment"  Check Deployment State  ${name}
    Run Keyword If  "${deployment_type}" == "daemonset"  Check Daemon Set State  ${name}

Check Daemon Set State With Prerequisite
    [Arguments]  ${name}  ${name-in-cr}  ${parentservice}=${None}
    ${status_check_object}=  Check that in CR service is presented  ${name-in-cr}  ${parentservice}
    ${flag}=  Run Keyword If  ${status_check_object}==True  Check Daemon Set State  ${name}
    RETURN  ${flag}

Check Daemon Set State
    [Arguments]  ${name}
    ${pods_in_namespace}=  Get Pods  ${namespace}
    ${pod_in_namespace}  Get Object In Namespace By Mask  ${pods_in_namespace}  ${name}
    ${daemon_set}=  Get Daemon Set  ${name}  ${namespace}
    Check Pod's List Is Equals  ${pod_in_namespace}  ${daemon_set.status.desired_number_scheduled}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Status Of Pods  ${pod_in_namespace}
    RETURN  ${flag}

Check Daemon Set And Deployment State For Cert Exporter
    [Arguments]  ${name}
    ${daemon_set}=  Get Daemon Set  ${name}  ${namespace}
    ${deployment}=  Get Deployment Entity  ${name}  ${namespace}
    ${pods_in_namespace}=  Get Pods  ${namespace}
    ${pod_in_namespace}  Get Object In Namespace By Mask  ${pods_in_namespace}  ${name}
    ${expected_sum_pods}=  Evaluate
    ...  ${daemon_set.status.desired_number_scheduled}+${deployment.spec.replicas}
    Check Pod's List Is Equals  ${pod_in_namespace}  ${expected_sum_pods}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Status Of Pods  ${pod_in_namespace}
    RETURN  ${flag}

Check Deployment State With Prerequisite
    [Arguments]  ${name}  ${name-in-cr}  ${parentservice}=${None}
    ${status_check_object}=  Check That In CR Service Is Presented  ${name-in-cr}  ${parentservice}
    ${flag}=  Run Keyword If  ${status_check_object}==True  Check Deployment State  ${name}
    RETURN  ${flag}

Check Deployment State
    [Arguments]  ${name}
    ${pods_in_namespace}=  Get Pods  ${namespace}
    ${pod_in_namespace}  Get Object In Namespace By Mask  ${pods_in_namespace}  ${name}
    ${deployment}=  Get Deployment Entity  ${name}  ${namespace}
    Check Pod's List Is Equals  ${pod_in_namespace}  ${deployment.spec.replicas}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Status Of Pods  ${pod_in_namespace}
    RETURN  ${flag}

Check Stateful Set State With Prerequisite
    [Arguments]  ${name}  ${name-in-cr}  ${parentservice}=${None}
    ${status_check_object}=  Check That In CR Service Is Presented  ${name-in-cr}  ${parentservice}
    ${flag}=  Run Keyword If  ${status_check_object}==True  Check Stateful Set State  ${name}
    RETURN  ${flag}

Check Stateful Set State
    [Arguments]  ${name}
    ${pods_in_namespace}=  Get Pods  ${namespace}
    ${pod_in_namespace}  Get Object In Namespace By Mask  ${pods_in_namespace}  ${name}
    ${stateful_set}=  Get Stateful Set  ${name}  ${namespace}
    Check Pod's List Is Equals  ${pod_in_namespace}  ${stateful_set.spec.replicas}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Status Of Pods  ${pod_in_namespace}
    RETURN  ${flag}

Check Prometheus Config Status
    ${resp}=  GET On Session  prometheussession  url=/api/v1/status/config
    Should Be Equal As Strings  ${resp.status_code}  200
    Should Contain  str(${resp.content})  "status":"success"

Check Prometheus Runtime Status
    ${resp}=  GET On Session  prometheussession  url=/api/v1/status/runtimeinfo
    Should Be Equal As Strings  ${resp.status_code}  200
    Should Contain  str(${resp.content})  "status":"success"

Check Prometheus Flags Status
    ${resp}=  GET On Session  prometheussession  url=/api/v1/status/flags
    Should Be Equal As Strings  ${resp.status_code}  200
    Should Contain  str(${resp.content})  "status":"success"

Get All Metrics From Api
    [Arguments]  ${session}
    ${response}=  GET On Session  ${session}  url=/api/v1/query?query=up
    Should Be Equal As Strings  ${response.status_code}  200
    RETURN  ${response}

Check Job Metrics Are Written
    [Arguments]  ${job}  ${metrics}
    ${job_metrics}=  Get Metrics By Job  ${metrics}  ${job}
    Run Keyword If  ${job_metrics}==False  Fail
    ...  Error! In ${metrics} metrics of ${job} don't exist
    ${status}=  Check Metrics Is Not Empty  ${job_metrics}
    Should Be Equal As Strings  ${status}  True

Check Kube State Metrics Target Metrics
    [Arguments]  ${kube_state_metrics_flag}  ${metrics}
    Run Keyword If  ${kube_state_metrics_flag}==True  Check Target is UP  ${kube-state-metrics}  ${all_active_targets}
    Run Keyword If  ${kube_state_metrics_flag}==True  Check Job Metrics Are Written  ${kube-state-metrics}  ${metrics}

Check Node Exporter Target Metrics
    [Arguments]  ${node_exporter_flag}  ${metrics}
    Run Keyword If  ${node_exporter_flag}==True  Check Target is UP  ${node-exporter}  ${all_active_targets}
    Run Keyword If  ${node_exporter_flag}==True  Check Job Metrics Are Written  ${node-exporter}  ${metrics}

Check Route/Ingress Status
    [Arguments]  ${name-in-cr}  ${name}  ${parentservice}=${None}  ${session}=external
    ${custom_resource}=  Get Custom Resource  monitoring.qubership.org/v1alpha1  PlatformMonitoring
    ...  ${namespace}  platformmonitoring
    ${external_url}=  Check Route Or Ingress  ${custom_resource}  ${name-in-cr}
    ...  ${namespace}-${name}  ${namespace}  ${parentservice}
    Preparation Session For External Service  ${external_url}  ${session}
    RETURN  True

Check Service Web UI Status Via External Url
    [Arguments]  ${url}  ${expected_string}  ${session}=external
    ${resp}=  GET On Session  ${session}  url=${url}
    Should Be Equal As Strings  ${resp.status_code}  200
    Should Contain  str(${resp.content})  ${expected_string}

Get All Active Targets
    [Arguments]  ${session}
    ${response}=  GET On Session  ${session}  url=/api/v1/targets?state=active
    Should Be Equal As Strings  ${response.status_code}  200
    RETURN  ${response}

Check Target Is UP
     [Arguments]  ${target_name}  ${all_active_targets}
     ${json_target}=  Get Prometheus Target  ${all_active_targets}  ${target_name}
     Run Keyword If  ${json_target}==False  Fail
     ...  Error! Target of ${target_name} doesn't exist
     ${status_check_object}=  Target State And Not Empty  ${json_target}  up
     Should Be Equal As Strings  ${status_check_object}  True

Check That In CR service Is Presented
     [Arguments]  ${name}  ${parentservice}
     ${custom_resource}=  Get Custom Resource  monitoring.qubership.org/v1alpha1  PlatformMonitoring  ${namespace}  platformmonitoring
     ${flag}=  Check CR Service Exists  ${custom_resource.get('spec')}  ${name}  ${parentservice}
     Skip If  ${flag} != True  Section ${name} is not presented in CR
     RETURN  ${flag}