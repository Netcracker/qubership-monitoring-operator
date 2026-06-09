*** Settings ***
Resource      ../smoke-test/keywords.robot
Resource      keywords.robot

*** Test Cases ***
Check Dashboards Availability In External Grafanas
    [Tags]    configurations-streamer
    Retrieve And Set Grafana Org Token    ${namespace}
    ${result}=    Check Dashboards Availability
    Log    ${result}

Create And Verify Simple Dashboard
    [Tags]    configurations-streamer
    Create And Verify Simple Dashboard    ${namespace}

Update Title And Verify Simple Dashboard
    [Tags]    configurations-streamer
    Update Title And Verify Dashboard    ${namespace}

Check Dashboards Availability On FTP
    [Tags]    configurations-streamer
    ${result}=    Compare FTP Dashboards
    Log    ${result}

Delete Test Dashboard
    [Tags]    configurations-streamer
    Remove Test Dashboard  ${namespace}  ${DASHBOARD_NAME}

Check Icinga Files On FTP
    [Tags]    configurations-streamer
    ${result}=    Check Alertmanager Files
    Log    ${result}

Check Prometheus Rules
    [Tags]    configurations-streamer
    ${result}=    Compare Prometheus Rules
    Log    ${result}

Resend Data Test
    [Tags]    configurations-streamer
    Wait Until Keyword Succeeds    2x    2s    Create API Session
    Wait Until Keyword Succeeds    2x    2s    Send Resend Request