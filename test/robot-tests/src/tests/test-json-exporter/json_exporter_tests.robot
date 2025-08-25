*** Settings ***
Library    RequestsLibrary
Library    Collections
Library    String
Library    BuiltIn
Resource   ../smoke-test/keywords.robot

Suite Setup     Preparation Operator Session
Suite Teardown  Close Session

*** Variables ***
${METRIC_SUFFIX}    _value
${RETRY_COUNT}      5
${RETRY_DELAY}      15s

*** Keywords ***
Get Metrics From ConfigMap
    ${configmap}=  Get Config Map  ${json-exporter}  ${namespace}
    ${config}=  Extract Config Map  ${configmap}  config.yml
    ${modules}=  Get From Dictionary    ${config}    modules
    ${tsdbstatus}=  Get From Dictionary    ${modules}    tsdbstatus
    ${raw_metrics}=  Get From Dictionary    ${tsdbstatus}    metrics
    ${metrics_with_suffix}=  Create List
    FOR  ${metric}  IN  @{raw_metrics}
        Append To List  ${metrics_with_suffix}  ${metric["name"]}${METRIC_SUFFIX}
    END
    RETURN  ${metrics_with_suffix}

Check Metrics In Prometheus Or VictoriaMetrics
    ${failed_metrics}=  Create List
    ${metrics}=  Get Metrics From ConfigMap
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

*** Test Cases ***
Check JSON Exporter Prometheus Targets
    [Tags]  full  smoke-test-prometheus  json-exporter
    ${all_active_targets}=  Get All Active Targets  prometheussession
    ${prometheus_metrics}=  Get All Metrics From Api  prometheussession
    Check Target Is UP  ${json-exporter}  ${all_active_targets}
    Check Job Metrics Are Written  ${json-exporter}  ${prometheus_metrics}
    
Check JSON Exporter Vmagent Targets
    [Tags]  full  smoke-test-vm  json-exporter
    ${all_active_targets}=  Get All Active Targets  vmagentsession
    ${vmagent_metrics}=  Get All Metrics From Api  vmsinglessession
    Check Target Is UP  ${json-exporter}  ${all_active_targets}
    Check Job Metrics Are Written  ${json-exporter}  ${vmagent_metrics}

Check JSON Exporter Metrics In Monitoring System
    [Tags]  full  json-exporter
    Check Metrics In Prometheus Or VictoriaMetrics