*** Settings ***
Library    Process
Library    RequestsLibrary
Library    Collections
Library    OperatingSystem
Library    String
Resource   ../smoke-test/keywords.robot

Suite Setup    Set Pushgateway Flag

*** Variables ***
${PUSHGATEWAY_HOST}    http://%{PUSHGATEWAY_HOST}
${METRICS_FILE}        ${EXECDIR}/tests/test-pushgateway/metrics_data.txt
${RETRY_COUNT}         6
${RETRY_DELAY}         15s

*** Keywords ***
Set Pushgateway Flag
    ${pushgateway_flag}=  Check Deployment State With Prerequisite  ${pushgateway}  ${pushgateway-in-cr}
    Set Suite Variable  ${pushgateway_flag}

Send Multiple Metrics
    ${command}  Set Variable  wget --post-file=${METRICS_FILE} ${PUSHGATEWAY_HOST}/metrics/job/multiple_job --no-check-certificate -O-
    ${result}  Run Process  ${command}  shell=True  timeout=10
    Should Be Equal As Strings  ${result.rc}  0

Check Metrics In Prometheus Or VictoriaMetrics
    ${failed_metrics}=  Create List
    ${metrics_data}=  Get File  ${METRICS_FILE}
    ${metrics}=  Split To Lines  ${metrics_data}
    FOR  ${metric}  IN  @{metrics}
        ${metric_name_list}=  Get Regexp Matches  ${metric}  ^([a-zA-Z_]+)
        ${metric_name}=  Set Variable  ${metric_name_list}[0]
        ${success}=  Set Variable  False
        FOR  ${i}  IN RANGE  ${RETRY_COUNT}
            ${response}=  Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
            ...  GET On Session  prometheussession  url=/api/v1/query?query=${metric_name}
            ...  ELSE  GET On Session  vmsinglessession  url=/api/v1/query?query=${metric_name}

            ${json_string}=  Evaluate  json.dumps(${response.json()}, ensure_ascii=False, sort_keys=True)  json
            ${expected_string}=  Set Variable  "__name__": "${metric_name}"

            ${status}=  Run Keyword And Return Status  Should Contain  ${json_string}  ${expected_string}
            Run Keyword If  ${status}  Set Test Variable  ${success}  True
            Exit For Loop If  ${success}

            Sleep  ${RETRY_DELAY}
        END
        Run Keyword If  not ${success}  Append To List  ${failed_metrics}  ${metric}
    END
    Run Keyword If  len(${failed_metrics}) > 0  Fail  Metrics not found: ${failed_metrics}

*** Test Cases ***
Push Metrics To Pushgateway
    [Tags]  full  pushgateway
    Run Keyword If  ${pushgateway_flag}  Send Multiple Metrics
    
Check Metrics In Monitoring System
    [Tags]  full  pushgateway
    [Setup]  Preparation Operator Session
    Run Keyword If  ${pushgateway_flag}  Check Metrics In Prometheus Or VictoriaMetrics