*** Settings ***

Resource        keywords.robot
Library         String

Suite Setup     Preparation Test Data
Suite Teardown  Cleanup Test Data


*** Keywords ***

Preparation Test Data
    Preparation Operator Session
    ${app_name}=  Get App Name From File  ${DEPLOYMENT}
    Set Suite Variable  ${app_name}
    Create Metric Test Application
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Pods Count Is  1  ${app_name}

Cleanup Test Data
    Delete All Sessions
    Delete Metric Test Application
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Pods Count Is  0  ${app_name}

*** Test Cases ***

Check Prometeus Target's Status, Url, Labels
    [Tags]  full  test-app
    ${target_object}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Target Of Test App Is Exist
    Check Target Of Test App Is UP And Has Labels  ${target_object}

Wait Until Metrics Of Test App Is Exist
    [Tags]  full  test-app
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Metrics Of Test App Is Exist

Check Available Metrics Of Test App
     [Tags]  full  test-app
     ${response}=  Check Metrics Of Test App Are Written
     Check Metrics Are Available And Not Empty  ${response}

Update Service Monitor Label To 'monitoring-test'
    [Tags]  full  test-app
    Update Service Monitor Label To  monitoring-test
    Check Service Monitor Label Updated To  monitoring-test

Target Of Test App Has Been Deleted And No Metrics Are Written
    [Tags]  full  test-app
    Add Selector To Cr  platformmonitoring  ${namespace}
    ${cr}=  Get Namespaced Custom Object Status  monitoring.netcracker.com  v1  ${namespace}  platformmonitorings  platformmonitoring
    Should Contain  str(${cr})  serviceMonitorSelector
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Target Of Test App Doesn't Exist
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check No Metrics Are Written
    [Teardown]  Delete Selector From Cr  platformmonitoring  ${namespace}

Return label: Update label in Custom Resource to 'monitoring'
    [Tags]  full  test-app
    Update Service Monitor Label To  monitoring
    Check Service Monitor Label Updated To  monitoring
