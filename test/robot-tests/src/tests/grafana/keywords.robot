*** Settings ***
Library            String
Library            json
Library            RequestsLibrary
Library            BuiltIn
Library            Collections
Library            PlatformLibrary                   managed_by_operator=true
Library            MonitoringLibrary
Resource           %{ROBOT_HOME}/tests/smoke-test/keywords.robot
Library            %{ROBOT_HOME}/lib/CheckJsonObject.py

*** Variables ***
${namespace}                %{NAMESPACE}
${grafana_host}             %{GRAFANA_HOST}
${FILES_PATH}               %{ROBOT_HOME}/source_files/dashboards

${PATH_TO_DASHBOARD}        ${FILES_PATH}/dashboard_for_create.yml
${PATH_TO_UPD_DASHBOARD}    ${FILES_PATH}/dashboard_for_update.yml
${RETRY_TIME}               5min
${RETRY_INTERVAL}           3s

*** Keywords ***
Initialize Grafana Library
    ${username}  ${password}=  Get Grafana Credentials From Secret
    Import Library   %{ROBOT_HOME}/lib/GrafanaApiLib.py
    ...              url=${grafana_host}
    ...              g_user=${username}
    ...              g_password=${password}

Get Grafana Credentials From Secret
    ${secret}=  Get Secret  grafana-admin-credentials  ${namespace}
    ${username_base64}=  Set Variable  ${secret.data["GF_SECURITY_ADMIN_USER"]}
    ${password_base64}=  Set Variable  ${secret.data["GF_SECURITY_ADMIN_PASSWORD"]}
    ${username}=  Decode Base64  ${username_base64}
    ${password}=  Decode Base64  ${password_base64}
    RETURN  ${username}  ${password}

Decode Base64
    [Arguments]  ${base64_string}
    ${decoded}=  Evaluate  base64.b64decode("""${base64_string}""").decode("utf-8")
    RETURN  ${decoded}

Attempt Login To Grafana
    [Arguments]  ${url}  ${username}  ${password}
    ${url}=  Set Variable  https://${url}
    Create Session  grafana_session  ${url}
    ${headers}=  Create Dictionary  Content-Type=application/json
    ${body}=  Create Dictionary  user=${username}  password=${password}
    ${response}=  POST On Session  grafana_session  /login
    ...  headers=${headers}
    ...  json=${body}
    ${login_status}=  Run Keyword If  '${response.status_code}' == '200'
    ...  Set Variable  SUCCESS
    ...  ELSE
    ...  Set Variable  FAIL
    Should Be Equal  ${login_status}  SUCCESS  Login to Grafana failed!
    RETURN  ${login_status}

Create Test Dashboard In Namespace
    [Arguments]  ${PATH_TO_DASHBOARD}
    ${body}=  Parse Yaml File  ${PATH_TO_DASHBOARD}
    ${created_dashboard}=  Create Dashboard In Namespace  ${namespace}  ${body}

Check That Dashboard Created Successfuly
    [Arguments]   ${dashboard_name}  ${namespace}
    ${object}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Resource Status Success in Cloud  ${namespace}  grafanadashboards  ${dashboard_name}
    ${status}=   Set Variable  ${object.get("status")}
    Run Keyword If  ${status} != None  Should Be Equal As Strings  ${status.get('message')}   success
    RETURN  ${object}

Check Resource Status Success in Cloud
    [Arguments]  ${namespace}  ${accepted_names:plural}  ${object_name}
    ${object}=  Get Namespaced Custom Object Status
    ...  integreatly.org  v1alpha1  ${namespace}  ${accepted_names:plural}  ${object_name}
    Should Not Be Equal  ${object}  ${NONE}
    RETURN  ${object}

Check That Grafana CR Includes Test Dashboard
    [Arguments]   ${dashboard_name}  ${namespace}
    ${object}=  Wait Until Keyword Succeeds  ${RETRY_TIME}  ${RETRY_INTERVAL}
    ...  Check Status Of Dashboards Includes Name, Namespace  ${namespace}  ${dashboard_name}
    ${status}=   Set Variable  ${object.get("status")}
    Run Keyword If  ${status} != None  Should Be Equal As Strings  ${status.get('message')}  success
    RETURN  ${object}

Check Status Of Dashboards Includes Name, Namespace
    [Arguments]  ${namespace}  ${name}
    ${grafana_CR_status}=   Check Resource Status Success in Cloud  ${namespace}  grafanas  grafana
    ${dashboard_status}=  Get Dashboard From Status  ${grafana_CR_status}  ${namespace}  ${name}
    Run Keyword If  ${dashboard_status} == None
    ...  Fail  Error! Dashboard with following name not found in namespace
    RETURN  ${dashboard_status}

Check Dashboard UID From Status Corresponds UID From File
    [Arguments]  ${status}  ${uid_from_file}
    ${uid_from_status}=  Get Dashboard Uid  ${status}
    Should Be Equal As Strings  ${uid_from_file}  ${uid_from_status}

Get Dashboard Uid
    [Arguments]  ${dashboard_status}
    ${uid}  Set Variable  ${dashboard_status.get('uid')}
    Set Suite Variable  ${uid}
    RETURN  ${uid}

Check Dashboard Is Appear In Grafana
    [Arguments]  ${uid}
    ${state}=  Run Keyword And Return Status  Find Dashboard  ${uid}
    Should Be True  ${state}

Delete Dashboard Via Cloud Rest
    [Arguments]  ${dashboard_name}
    ${delete_status}=  Delete Dashboard In Namespace  ${namespace}  ${dashboard_name}
    Should Be Equal As Strings  ${delete_status.get('status')}  Success

Check Dashboard Is Deleted In Grafana
    [Arguments]  ${uid}
    ${state}=  Run Keyword And Return Status  Find Dashboard  ${uid}
    Should Be Equal As Strings  ${state}  False

Prepare Data For Update Dashboard
    [Arguments]  ${PATH_TO_UPD_DASHBOARD}  ${namespace}  ${dashboard_name}
    ${dashboard_from_file}=  Parse Yaml File  ${PATH_TO_UPD_DASHBOARD}
    ${dashboard}=  Check Existing Dashboard And Create It If Not Found  ${namespace}  ${dashboard_name}
    # For update dashboard custom resource it need to update dashboard ResourceVersion from file to ResourceVersion from dashboard in cloud
    ${resourceVersion}=  Get ResourceVersion From Dashboard  ${dashboard}
    ${updated_dashboard}=  Update Dashboard Parameter  ${dashboard_from_file}  resourceVersion  ${resourceVersion}
    RETURN  ${updated_dashboard}

Check Existing Dashboard And Create It If Not Found
    [Arguments]  ${namespace}  ${name}
    ${status}  ${dashboard}=  Run Keyword And Ignore Error
    ...  Get Dashboard In Namespace  ${namespace}  ${name}
    ${tmp_dashboard}=  Set Variable  ${dashboard}
    ${dashboard}=  Run Keyword If  '${status}'=='FAIL'
    ...  Create Test Dashboard In Namespace  ${PATH_TO_DASHBOARD}
    ...  ELSE  Set Variable  ${tmp_dashboard}
    RETURN  ${dashboard}

Get ResourceVersion From Dashboard
    [Arguments]  ${dashboard}
    ${resourceVersion}=  Set Variable  ${dashboard.get('metadata').get('resourceVersion')}
    RETURN  ${resourceVersion}

Get Dashboard And Check It's Updated In Cloud
    [Arguments]  ${namespace}  ${name}  ${uid}  ${dashboard_from_file}
    ${dashboard}=  Get Dashboard In Namespace  ${namespace}  ${name}
    ${dashboard_json}=  Evaluate   ${dashboard.get('spec').get('json')}
    ${dashboard_json}=  Convert To Json  ${dashboard_json}
    ${dashboard_json_from_file}=  Evaluate  ${dashboard_from_file.get('spec').get('json')}
    ${dashboard_json_from_file}=  Convert To Json  ${dashboard_json_from_file}
    Should Be Equal  ${dashboard_json.get('title')}  updated-robotframework-dashboard
    Should Be Equal  ${dashboard_json.get('uid')}  ${uid}
    ${result}=  Compare Two Jsons  ${dashboard_json}  ${dashboard_json_from_file}
    Should Be Equal  ${result}  ${TRUE}

Get Dashboard and Check it's Updated in Grafana
    [Arguments]  ${uid}  ${dashboard_from_file}
    ${dashboard_from_Grafana}=  Find Dashboard  ${uid}
    Should Be Equal  ${dashboard_from_Grafana.get('dashboard').get('title')}  updated-robotframework-dashboard
    Should Be Equal  ${dashboard_from_Grafana.get('dashboard').get('uid')}  ${uid}
    ${dashboard_json_from_file}=  Evaluate  ${dashboard_from_file.get('spec').get('json')}
    ${dashboard_json_from_file}=  Convert To Json  ${dashboard_json_from_file}
    ${result}=  Compare Two Jsons  ${dashboard_from_Grafana.get('dashboard')}  ${dashboard_json_from_file}
    Should Be Equal  ${result}  ${TRUE}
