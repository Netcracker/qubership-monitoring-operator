*** Settings ***
Resource  keywords.robot
Library   String

Suite Setup  Initialize And Prepare Test Data
Suite Teardown  Cleanup Test Data

*** Keywords ***
Cleanup Test Data
    ${status}  ${message}=  Run Keyword And Ignore Error  Find Dashboard  ${uid}
    Run Keyword If  '${status}'=='PASS'  Delete Dashboard Via Cloud Rest

Initialize And Prepare Test Data
    Initialize Grafana Library
    ${uid_from_file}=  Get Dashboard Value From File   ${PATH_TO_DASHBOARD}  uid
    ${dashboard_name}=  Get Dashboard Value From File   ${PATH_TO_DASHBOARD}  title
    Set Suite Variable  ${uid_from_file}
    Set Suite Variable  ${dashboard_name}

*** Test Cases ***
Verify Login To Grafana
    [Tags]  full  grafana
    ${username}  ${password}=  Get Grafana Credentials From Secret
    Log to console  username ${username}
    Log to console  password ${password}
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Attempt Login To Grafana  ${grafana_host}  ${username}  ${password}

Create Object Grafana Dashboard
    [Tags]  full  grafana
    Create Test Dashboard In namespace  ${PATH_TO_DASHBOARD}
    Check That Dashboard Created Successfuly  ${dashboard_name}  ${namespace}
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Dashboard Is Appear In Grafana  ${uid_from_file}

Check Dashboard Is Appear In Grafana CR, Includes All Fields
    [Tags]  full  grafana_checkCR
    ${status_from_grafana_CR}=  Check That Grafana CR Includes Test Dashboard  ${dashboard_name}  ${namespace}
    Check Dashboard UID From Status Corresponds UID From File  ${status_from_grafana_CR}  ${uid_from_file}
    Check Dashboard Is Appear In Grafana  ${uid_from_file}

Update GrafanaDashboard And Check It Updated
    [Tags]  full  grafana
    ${updated_dashboard}=  Prepare Data For Update Dashboard
    ...  ${PATH_TO_UPD_DASHBOARD}  ${namespace}  ${dashboard_name}
    Replace Dashboard In Namespace  ${namespace}  ${dashboard_name}  ${updated_dashboard}
    ${state}=  Run Keyword And Return Status  Get Dashboard and Check It's Updated In Cloud
    ...  ${namespace}  ${dashboard_name}  ${uid_from_file}  ${updated_dashboard}
    Run Keyword If  '${state}'=='FAIL'  Fail
    ...  The 'GrafanaDashboard' Custom Resource is not updated in the Cloud!
    ${state}=  Run Keyword And Return Status
    ...  Get Dashboard and Check it's Updated in Grafana  ${uid_from_file}  ${updated_dashboard}
    Run Keyword If  '${state}'=='FAIL'  Fail
    ...  The 'GrafanaDashboard' is not updated in the Grafana!

Delete Dashboard And Check It Deleted In Grafana
    [Tags]  full  grafana
    Delete Dashboard Via Cloud Rest  ${dashboard_name}
    Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Dashboard is deleted in Grafana  ${uid_from_file}
