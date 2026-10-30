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
    ${use_url}=  Set Variable  ${external_url}
    ${has_scheme}=  Run Keyword And Return Status  Should Start With  ${external_url}  http
    IF  not ${has_scheme}
        ${use_url}=  Determine Protocol  ${external_url}  ${auth}
    END
    Create Session  ${session}  ${use_url}  auth=${auth}  verify=${False}

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
    Get Daemon Set  ${name}  ${namespace}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Daemon Set Is Rolled Out And Pods Are Running  ${name}
    RETURN  ${flag}

Check Daemon Set And Deployment State For Cert Exporter
    [Arguments]  ${name}
    Get Daemon Set  ${name}  ${namespace}
    Get Deployment Entity  ${name}  ${namespace}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Cert Exporter Is Rolled Out And Pods Are Running  ${name}
    RETURN  ${flag}

Check Deployment State With Prerequisite
    [Arguments]  ${name}  ${name-in-cr}  ${parentservice}=${None}
    ${status_check_object}=  Check That In CR Service Is Presented  ${name-in-cr}  ${parentservice}
    ${flag}=  Run Keyword If  ${status_check_object}==True  Check Deployment State  ${name}
    RETURN  ${flag}

# The workload is fetched once up front so that a missing one fails immediately; only its rollout
# and pods are waited for, within a single RETRY_TIME budget.
Check Deployment State
    [Arguments]  ${name}
    Get Deployment Entity  ${name}  ${namespace}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Deployment Is Rolled Out And Pods Are Running  ${name}
    RETURN  ${flag}

Check Stateful Set State With Prerequisite
    [Arguments]  ${name}  ${name-in-cr}  ${parentservice}=${None}
    ${status_check_object}=  Check That In CR Service Is Presented  ${name-in-cr}  ${parentservice}
    ${flag}=  Run Keyword If  ${status_check_object}==True  Check Stateful Set State  ${name}
    RETURN  ${flag}

Check Stateful Set State
    [Arguments]  ${name}
    Get Stateful Set  ${name}  ${namespace}
    ${flag}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Stateful Set Is Rolled Out And Pods Are Running  ${name}
    RETURN  ${flag}

# A rollout is complete when the controller has observed the current generation and every replica is
# updated and available, the same conditions `kubectl rollout status` waits for. Checking pods before
# that point races with the rolling update: the old pod is still Running while it terminates.
# These keywords re-read the workload and the pod list on every call so that they can be retried.
Deployment Is Rolled Out And Pods Are Running
    [Arguments]  ${name}
    ${deployment}=  Get Deployment Entity  ${name}  ${namespace}
    Rollout Is Complete  ${deployment}
    ${flag}=  Check Workload Pods Are Running  ${name}  ${deployment.spec.replicas}
    RETURN  ${flag}

Stateful Set Is Rolled Out And Pods Are Running
    [Arguments]  ${name}
    ${stateful_set}=  Get Stateful Set  ${name}  ${namespace}
    Rollout Is Complete  ${stateful_set}
    ${flag}=  Check Workload Pods Are Running  ${name}  ${stateful_set.spec.replicas}
    RETURN  ${flag}

Daemon Set Is Rolled Out And Pods Are Running
    [Arguments]  ${name}
    ${daemon_set}=  Get Daemon Set  ${name}  ${namespace}
    Rollout Is Complete  ${daemon_set}
    ${flag}=  Check Workload Pods Are Running  ${name}  ${daemon_set.status.desired_number_scheduled}
    RETURN  ${flag}

Cert Exporter Is Rolled Out And Pods Are Running
    [Arguments]  ${name}
    ${daemon_set}=  Get Daemon Set  ${name}  ${namespace}
    ${deployment}=  Get Deployment Entity  ${name}  ${namespace}
    Rollout Is Complete  ${daemon_set}
    Rollout Is Complete  ${deployment}
    ${expected_sum_pods}=  Evaluate
    ...  ${daemon_set.status.desired_number_scheduled}+${deployment.spec.replicas}
    ${flag}=  Check Workload Pods Are Running  ${name}  ${expected_sum_pods}
    RETURN  ${flag}

Rollout Is Complete
    [Arguments]  ${workload}
    ${reason}=  Describe Incomplete Rollout  ${workload}
    Should Be Empty  ${reason}  ${reason}

Check Workload Pods Are Running
    [Arguments]  ${name}  ${expected_count}
    ${pods_in_namespace}=  Get Pods  ${namespace}
    ${pods}=  Get Object In Namespace By Mask  ${pods_in_namespace}  ${name}
    ${pods}=  Exclude Terminating Pods  ${pods}
    Check Pod's List Is Equals  ${pods}  ${expected_count}
    ${flag}=  Check Status Of Pods  ${pods}
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
    ${custom_resource}=  Get Custom Resource  monitoring.netcracker.com/v1  PlatformMonitoring
    ...  ${namespace}  platformmonitoring
    # Retry while wildcard Ingress LB hostname is still pending; API errors also
    # surface here instead of being treated as a soft skip.
    ${external_url}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Route Or Ingress  ${custom_resource}  ${name-in-cr}
    ...  ${namespace}-${name}  ${namespace}  ${parentservice}
    Run Keyword If  '${external_url}' == 'None'  Return From Keyword  False
    Run Keyword If  '${external_url}' == ''  Return From Keyword  False
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
     ${custom_resource}=  Get Custom Resource  monitoring.netcracker.com/v1  PlatformMonitoring  ${namespace}  platformmonitoring
     ${flag}=  Check CR Service Exists  ${custom_resource.get('spec')}  ${name}  ${parentservice}
     Skip If  ${flag} != True  Section ${name} is not presented in CR
     RETURN  ${flag}

Check Target And Metrics Are Ready
    [Arguments]  ${target}  ${all_active_targets}  ${metrics}
    Check Target Is UP  ${target}  ${all_active_targets}
    Check Job Metrics Are Written  ${target}  ${metrics}

Check Target And Metrics Are Ready Refreshing
    [Arguments]  ${target}  ${targets_session}  ${metrics_session}
    ${all_active_targets}=  Get All Active Targets  ${targets_session}
    ${metrics}=             Get All Metrics From Api  ${metrics_session}
    ${json_target}=  Get Prometheus Target  ${all_active_targets}  ${target}
    Run Keyword If  ${json_target}==${False}
    ...  Fail  Target ${target} missing from ${targets_session} active targets
    ${is_up}=  Target State And Not Empty  ${json_target}  up
    Run Keyword If  not ${is_up}
    ...  Log  Target ${target} not UP yet (session=${targets_session}): ${json_target}  level=WARN
    Should Be Equal As Strings  ${is_up}  True
    ${job_metrics}=  Get Metrics By Job  ${metrics}  ${target}
    Run Keyword If  ${job_metrics}==${False}
    ...  Log  Job ${target} metrics not present in ${metrics_session} yet (up query result has no matching job)  level=WARN
    Run Keyword If  ${job_metrics}==${False}
    ...  Fail  Error! In ${metrics_session} metrics of ${target} don't exist
    ${status}=  Check Metrics Is Not Empty  ${job_metrics}
    Should Be Equal As Strings  ${status}  True

Check Vmagent Target Metrics With Retry
    [Arguments]  ${target}  ${all_active_targets}  ${metrics}
    ${flag}=  Wait Until Keyword Succeeds
    ...  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Target And Metrics Are Ready Refreshing
    ...  ${target}  vmagentsession  vmsinglessession
    RETURN  ${flag}

Check Kube State Metrics Are Ready
    [Arguments]  ${metrics}
    Check Target Is UP  ${kube-state-metrics}  ${all_active_targets}
    Check Job Metrics Are Written  ${kube-state-metrics}  ${metrics}

Check Kube State Metrics With Retry
    [Arguments]  ${kube_state_metrics_flag}  ${metrics}
    ${flag}=  Run Keyword If  ${kube_state_metrics_flag}==True
    ...  Wait Until Keyword Succeeds
    ...  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Kube State Metrics Are Ready
    ...  ${metrics}
    RETURN  ${flag}
