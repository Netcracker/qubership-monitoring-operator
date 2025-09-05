*** Settings ***
Resource           ../keywords.robot

*** Variables ***
${namespace}                        %{NAMESPACE}
${OPERATOR}                         %{OPERATOR}
${FILES_PATH}                       %{ROBOT_HOME}/source_files/test_adapter_app

${DEPLOYMENT}                       ${FILES_PATH}/deployment.yaml
${SERVICE}                          ${FILES_PATH}/service.yaml
${SERVICE_MONITOR}                  ${FILES_PATH}/servicemonitor.yaml
${METRIC_RULE}                      ${FILES_PATH}/custom-scale-metric-rule.yaml
${HORIZONTAL_POD_AUTOSCALER}        ${FILES_PATH}

${RETRY_TIME}                       8min
${RETRY_INTERVAL}                   1s
${autoscaler_url}                   http://autoscaling-example-service:8080

*** Keywords ***
Get Value Of prometheus_example_app_load In Prometheus
    ${response}=  Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
    ...  GET On Session  prometheussession  url=/api/v1/query?query=prometheus_example_app_load
    ...  ELSE  GET On Session  vmsinglessession  url=/api/v1/query?query=prometheus_example_app_load
    Should Be Equal As Strings  ${response.status_code}  200
    RETURN  ${response}

Check prometheus-adapter Changed Metric To
    [Arguments]  ${expected_value}
    ${response}=  Get Value Of prometheus_example_app_load In Prometheus
    ${metric_value}=  Get Value From Prometheus Query  ${response}
    Convert To Integer  ${expected_value}
    Should Be Equal As Integers  ${metric_value}  ${expected_value}

Preparation Sessions
    Preparation Operator Session
    Create Session  autoscalersession  ${autoscaler_url}

Create Adapter Test Application
    ${new_deployment}=  Add Security Context To Deployment  ${DEPLOYMENT}  ${namespace}
    Create Deployment Entity  ${new_deployment}  ${namespace}
    Create Service From File  ${SERVICE}  ${namespace}
    Create Service Monitor  ${SERVICE_MONITOR}
    Create Custom Metric Rule  ${METRIC_RULE}
    Create Horizontal Pod Autoscaler  ${HORIZONTAL_POD_AUTOSCALER}

Delete AutoScaler Test Application
    Delete Deployment Entity  ${app_name}  ${namespace}
    Delete Service  ${app_name}  ${namespace}
    Delete Service Monitor  ${app_name}
    Delete Custom Metric Rule  ${app_name}-custom-metric-rule
    Delete Horizontal Pod Autoscaler  ${app_name}
