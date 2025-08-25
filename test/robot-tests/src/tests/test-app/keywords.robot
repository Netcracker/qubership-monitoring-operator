# Useful links
# API kuber with python https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/AppsV1Api.md

*** Settings ***
Resource           ../keywords.robot

*** Variables ***
${namespace}                %{NAMESPACE}
${FILES_PATH}               %{ROBOT_HOME}/source_files/test_metric_app

${DEPLOYMENT}               ${FILES_PATH}/deployment.yml
${SERVICE}                  ${FILES_PATH}/service.yml
${SERVICEMONITOR}           ${FILES_PATH}/servicemonitor.yml
${job}                      ${service_name}
${target_testapp_name}      ${service_name}
${service_name}             prometheus-example-app

${RETRY_TIME}               8min
${RETRY_INTERVAL}           1s

*** Keywords ***
Close Session
    Delete All Sessions
    Delete Metric Test Application

Create Metric Test Application
    ${new_deployment}=  Add Security Context To Deployment  ${DEPLOYMENT}  ${namespace}
    Create Deployment Entity  ${new_deployment}  ${namespace}
    Create Service From File  ${SERVICE}  ${namespace}
    Create Service Monitor  ${SERVICE_MONITOR}

Delete Metric Test Application
    Delete Deployment Entity  ${app_name}  ${namespace}
    Delete Service  ${app_name}  ${namespace}
    Delete Service Monitor  ${app_name}

Get All Active Targets
    [Arguments]  ${session}
    ${response}=  GET On Session  ${session}  url=/api/v1/targets?state=active
    Should Be Equal As Strings  ${response.status_code}  200
    RETURN  ${response}

Check Target Of Test App Is Exist
    ${json_response}=  Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
    ...  Get All Active Targets  prometheussession
    ...  ELSE  Get All Active Targets  vmagentsession
     ${test_target}=  Get Prometheus Target  ${json_response}  ${target_testapp_name}
     Should Not Contain  ${test_target}  False
     RETURN  ${test_target}

Check Metrics Of Test App Is Exist
     ${prometheus_metrics}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
     ...  Check Metrics Of Test App Are Written
     Should Not Be Empty  ${prometheus_metrics.get('metric')}

Check Target Of Test App Is UP And Has Labels
     [Arguments]  ${json_target}
     Run Keyword If  ${json_target}==False  Fail
     ...  Error! Target Of Test Application Doesn't Exist
     ${status_check_object}=  Target State And Not Empty  ${json_target}  up
     ${check_result}=  Run Keyword And Return If  ${status_check_object} != 'False'
     ...  Check Labels In Target  ${json_target}  ${namespace}  ${service_name}  ${job}
     Should Be Equal As Strings  ${check_result}  True

Get All Metrics From Api
    ${response}=  Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
    ...  GET On Session  prometheussession  url=/api/v1/query?query=up
    ...  ELSE  GET On Session  vmsinglessession  url=/api/v1/query?query=up
    Should Be Equal As Strings  ${response.status_code}  200
    RETURN  ${response}

Check Metrics Of Test App Are Written
    ${all_metrics}=  Get All Metrics From Api
    ${prometheus_metrics}=  Get Metrics Of Test App  ${all_metrics}  ${job}  ${namespace}
    Should Not Be Equal As Strings  ${prometheus_metrics}  False
    RETURN  ${prometheus_metrics}

Check Metrics Are Available And Not Empty
    [Arguments]  ${prometheus_metrics}
    Run Keyword If  ${prometheus_metrics}==False  Fail
    ...  Error! Metrics of test application doesn't exist!
    ${status}=  Check Metrics Is Not Empty  ${prometheus_metrics}
    Should Be Equal As Strings  ${status}  True

Create Body for Update ServiceMonitor Label
    [Arguments]  ${label}
    ${data}=  Create Dictionary  app.kubernetes.io/component=${label}
    ${labels}=  Create Dictionary  labels=${data}
    ${metadata}=  Create Dictionary  metadata=${labels}
    RETURN  ${metadata}

Update Service Monitor Label To
    [Arguments]  ${label}
    ${update_data}=  Create Body for Update ServiceMonitor Label  ${label}
    Patch Service Monitor  ${service_name}  ${update_data}

Check Service Monitor Label Updated To
    [Arguments]  ${label}
    ${service_monitor}=  Get Service Monitor  ${service_name}
    ${status}=  Check CR Label Updated  ${service_monitor.get('metadata')}  ${label}
    Run Keyword If  ${status}==False  Fail
    ...  Error! Service Monitor labels of Test Application is not updated!
    Should Be Equal As Strings  ${status}  True

Target Of Test App Doesn't Exist
    ${all_active_targets}=  Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
    ...  Get All Active Targets  prometheussession
    ...  ELSE  Get All Active Targets  vmagentsession
    ${status}=  Get Prometheus Target  ${all_active_targets}  ${target_testapp_name}
    Should Be Equal As Strings  ${status}  False

Check No Metrics Are Written
    ${all_metrics}=  Get All Metrics From Api
    ${prometheus_metrics}=  Get Metrics Of Test App  ${all_metrics}  ${job}  ${namespace}
    Should Be Equal As Strings  ${prometheus_metrics}  False
