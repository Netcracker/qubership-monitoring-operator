*** Settings ***
Resource        keywords.robot
Library         String

Suite Setup     Preparation Test Data
Suite Teardown  Cleanup Test Data

*** Keywords ***
Preparation Test Data
    Preparation Sessions
    ${app_name}=  Get App Name From File  ${DEPLOYMENT}
    Set Suite Variable  ${app_name}
    Create Adapter Test Application
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Pods Count Is  1  ${app_name}

Cleanup Test Data
    Delete All Sessions
    Delete AutoScaler Test Application
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Pods Count Is  0  ${app_name}

Set prometheus_example_app_load To 5
    ${resp}=  GET On Session   autoscalersession  url=/load  timeout=10
    log to console  resp: ${resp.content}
    Should Be Equal As Strings  ${resp.status_code}  200
    Should Be Equal As Strings  ${resp.content}  Value of prometheus_example_app_load is set to 5

Set prometheus_example_app_load To 1
    ${resp}=  GET On Session   autoscalersession  url=/unload  timeout=10
    Should Be Equal As Strings  ${resp.status_code}  200
    Should Be Equal As Strings  ${resp.content}  Value of prometheus_example_app_load is set to 1

*** Test Cases ***
Increase Autoscaler Replicas To 3 And Check Metric Is 5
    [Tags]  full  prometheus-adapter
    Wait Until Keyword Succeeds  ${RETRY_TIME}  10s
    ...  Set prometheus_example_app_load To 5
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Pods Count Is  3  ${app_name}
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check prometheus-adapter Changed Metric To  5

Decrease Autoscaler Replicas To 2 And Check Metric Is 1
    [Tags]  full  prometheus-adapter
    Wait Until Keyword Succeeds  ${RETRY_TIME}  10s
    ...  Set prometheus_example_app_load To 1
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Pods Count Is  2  ${app_name}
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check prometheus-adapter Changed Metric To  1
