*** Settings ***
Library    RequestsLibrary
Library    Collections
Library    BuiltIn
Library    %{ROBOT_HOME}/lib/ConfigurationsStreamerLib.py

*** Variables ***
${GRAFANA_BASIC_USER}    admin
${GRAFANA_BASIC_PASS}    admin
${IP_GRAFANA_CREDS}      %{URL_GRAFANA_CREDS}
${IP_GRAFANA_TOKEN}      %{URL_GRAFANA_TOKEN}
${SECRET_NAME}           grafana-api-token
${DASHBOARD_NAME}        simple-dashboard
${PATH_TO_DASHBOARD}     %{ROBOT_HOME}/source_files/dashboards/simple-dashboard.yml
${DASHBOARD_TITLE}       Simple Dashboard
${DASHBOARD_TITLE_NEW}   Simple Dashboard Changed
${IP_FTP_SERVER}         %{FTP_SERVER}
${SSH_KEY}               %{SSH_KEY}
${FTP_VM_USER}           %{VM_USER}
${RESEND_BASE_URL}       http://configurations-streamer.${namespace}.svc.cluster.local:8282/api/v1/send
&{RESEND_AUTH_HEADER}    Authorization=Basic YWRtaW46YWRtaW4=

*** Keywords ***
Get Encoded Grafana Token
    [Arguments]    ${secret_name}    ${namespace}
    ${secret}=    Get Secret    ${secret_name}    ${namespace}
    ${encoded_token}=    Set Variable    ${secret.data["requestToken"]}
    RETURN    ${encoded_token}

Decode Base64 Token
    [Arguments]    ${encoded_token}
    ${decoded_token}=    Evaluate    __import__('base64').b64decode("""${encoded_token}""").decode('utf-8')    modules=base64
    RETURN    ${decoded_token}

Get Grafana Org Token
    [Arguments]    ${secret_name}    ${namespace}
    ${encoded_token}=    Get Encoded Grafana Token    ${secret_name}    ${namespace}
    ${org_token}=    Decode Base64 Token    ${encoded_token}
    RETURN    ${org_token}

Retrieve And Set Grafana Org Token
    [Arguments]    ${namespace}
    ${org_token}=    Get Grafana Org Token    ${SECRET_NAME}    ${namespace}
    Set Grafana Org Token    ${org_token}

Create And Verify Simple Dashboard
    [Arguments]    ${namespace}
    Retrieve And Set Grafana Org Token    ${namespace}
    ${body}=  Parse Yaml File  ${PATH_TO_DASHBOARD}
    Create Test Dashboard  ${namespace}  ${body}
    Sleep    30s
    Check Dashboard Existence  ${DASHBOARD_TITLE}

Update Title And Verify Dashboard
    [Arguments]    ${namespace}
    Retrieve And Set Grafana Org Token    ${namespace}
    Update Test Dashboard  ${namespace}  ${DASHBOARD_NAME}  ${DASHBOARD_TITLE_NEW}
    Sleep    30s
    Check Dashboard Existence  ${DASHBOARD_TITLE_NEW}
    
Create API Session
    Create Session    streamer    ${RESEND_BASE_URL}    headers=&{RESEND_AUTH_HEADER}    timeout=5

Send Resend Request
    ${response}    POST On Session    streamer    ${EMPTY}    expected_status=200
    ${json_response}    Set Variable    ${response.json()}
    Should Be Equal As Strings    ${json_response}[message]    Data resent