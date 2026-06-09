*** Settings ***
Library    Collections
Library    OperatingSystem
Library    BuiltIn
Resource   ../smoke-test/keywords.robot

Suite Setup    Custom Suite Setup
Suite Teardown  Close Session

*** Variables ***
${ssh_collector}           %{SSH_COLLECTOR}
${http_collector}          %{HTTP_COLLECTOR}
${RETRY_COUNT}             5
${RETRY_DELAY}             10s

*** Keywords ***
Custom Suite Setup
    Preparation Operator Session
    Set Version Exporter Flag
    Parse Version Exporter ConfigMap
    Wait Before Testing

Wait Before Testing
    Log To Console  Waiting for startup (completion of ConfigMap scanning) version-exporter
    Sleep  200s

Set Version Exporter Flag
    ${version-exporter_flag}=  Check Deployment State  ${version-exporter}
    Set Suite Variable  ${version-exporter_flag}

Parse Version Exporter ConfigMap
    ${configmap}=  Get Config Map  ${version-exporter}  ${namespace}
    ${config}=  Extract Config Map  ${configmap}  exporterConfig.yaml
    Set Suite Variable  ${config}

Get Metrics From ConfigMap
    [Arguments]  ${collector_type}
    ${collector}=  Get From Dictionary  ${config}  ${collector_type}  default=${None}
    ${metrics}=  Create List
    Run Keyword If  '${collector_type}' == 'configmap_collector'  Get ConfigMap Metrics  ${collector}  ${metrics}
    Run Keyword If  '${collector_type}' != 'configmap_collector'  Get Standard Metrics  ${collector}  ${metrics}
    RETURN  ${metrics}

Get ConfigMap Metrics
    [Arguments]  ${collector}  ${metrics}
    ${defaults}=  Get From Dictionary  ${collector}  defaults
    ${resources}=  Get From Dictionary  ${collector}  resources  default=[]
    ${metric_name}=  Get From Dictionary  ${defaults}  metricName
    Append To List  ${metrics}  ${metric_name}
    FOR  ${resource}  IN  @{resources}
        ${metric_name}=  Get From Dictionary  ${resource}  metricName  default=${EMPTY}
        Run Keyword If  '${metric_name}' != '${EMPTY}'  Append To List  ${metrics}  ${metric_name}
    END

Get Standard Metrics
    [Arguments]  ${collector}  ${metrics}
    ${connections}=  Get From Dictionary  ${collector}  connections  default=[]
    FOR  ${conn}  IN  @{connections}
        ${requests}=  Get From Dictionary  ${conn}  requests  default=[]
        FOR  ${resource}  IN  @{requests}
            ${resource_dict}=  Convert To Dictionary  ${resource}
            ${metric_name}=  Get From Dictionary  ${resource_dict}  metricName
            Append To List  ${metrics}  ${metric_name}
        END
    END

Check Metrics In Prometheus Or VictoriaMetrics
    [Arguments]  ${metrics}
    ${failed_metrics}=  Create List
    FOR  ${metric}  IN  @{metrics}
        ${success}=  Set Variable  False
        FOR  ${i}  IN RANGE  ${RETRY_COUNT}
            ${response}=  Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
            ...  GET On Session  prometheussession  url=/api/v1/query?query=${metric}
            ...  ELSE  GET On Session  vmsinglessession  url=/api/v1/query?query=${metric}

            ${json_string}=  Evaluate  json.dumps(${response.json()}, ensure_ascii=False, sort_keys=True)  json
            ${expected_string}=  Set Variable  "__name__": "${metric}"

            ${status}=  Run Keyword And Return Status  Should Contain  ${json_string}  ${expected_string}
            Run Keyword If  ${status}  Set Test Variable  ${success}  True
            Exit For Loop If  ${success}

            Sleep  ${RETRY_DELAY}
        END
        Run Keyword If  not ${success}  Append To List  ${failed_metrics}  ${metric}
    END
    Run Keyword If  len(${failed_metrics}) > 0  Fail  Metrics not found: ${failed_metrics}

Validate SSH Metrics
    ${metrics}=  Get Metrics From ConfigMap  ssh_collector
    Check Metrics In Prometheus Or VictoriaMetrics  ${metrics}

Validate ConfigMap Metrics
    ${metrics}=  Get Metrics From ConfigMap  configmap_collector
    Check Metrics In Prometheus Or VictoriaMetrics  ${metrics}

Validate HTTP Metrics
    ${metrics}=  Get Metrics From ConfigMap  http_collector
    Check Metrics In Prometheus Or VictoriaMetrics  ${metrics}

*** Test Cases ***
Test ConfigMap Collector Metrics
    [Tags]  full  version-exporter
    Run Keyword If  '${version-exporter_flag}' == 'True'  Validate ConfigMap Metrics
    
Test HTTP Collector Metrics
    [Tags]  full  version-exporter
    Run Keyword If  '${http_collector}' != 'true'  Skip  HTTP Collector is disabled
    Validate HTTP Metrics
    
Test SSH Collector Metrics
    [Tags]  full  version-exporter
    Run Keyword If  '${ssh_collector}' != 'true'  Skip  SSH Collector is disabled
    Validate SSH Metrics