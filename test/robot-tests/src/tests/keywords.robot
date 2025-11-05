# https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/AppsV1Api.md
*** Settings ***
Library            String
Library            json
Library            RequestsLibrary
Library            BuiltIn
Library            Collections
Library            PlatformLibrary                      managed_by_operator=true
Library            %{ROBOT_HOME}/lib/CheckJsonObject.py
Library            %{ROBOT_HOME}/lib/TestAppsLib.py     managed_by_operator=true

*** Variables ***
${prometheus_host}          prometheus-operated:9090
${vmsingle_host}            vmsingle-k8s:8429
${vmagent_host}             vmagent-k8s:8429
${vmauth_host}              vmauth-k8s:8427
${vmuser}                   vmuser-k8s
${vmauth-in-cr}             vmAuth
${vmauth}                   False
${OPERATOR}                 %{OPERATOR}

*** Keywords ***
Determine Protocol
    [Arguments]  ${host}  ${auth}=None
    ${https_url}=  Set Variable  https://${host}
    ${http_url}=   Set Variable  http://${host}
    ${https_response}=  Run Keyword And Return Status  Check URL Accessibility  ${https_url}  ${auth}
    IF  ${https_response}
        Log To Console  Using HTTPS: ${https_url}
        ${final_url}=  Set Variable  ${https_url}
    ELSE
        Log To Console  Using HTTP: ${http_url}
        ${final_url}=  Set Variable  ${http_url}
    END
    RETURN  ${final_url}

Check URL Accessibility
    [Arguments]  ${url}  ${auth}=None
    Evaluate  __import__("logging").getLogger("urllib3").setLevel(40)
    Create Session  temp_session  ${url}  auth=${auth}
    ${response}=  GET On Session  temp_session  /
    Delete All Sessions
    Should Be Equal As Integers  ${response.status_code}  200
    RETURN  True

Get Creadentials From Secret
    ${secret}=  Get Secret  ${vmuser}  ${namespace}
    ${pass}=  Get Pass From Secret  ${secret}
    ${username}=  Get Username From Secret  ${secret}
    ${auth}=  Create List  ${username}  ${pass}
    RETURN  ${auth}

Get App Name From File
    [Arguments]  ${FILE_PATH}
    ${body}=  Parse Yaml File  ${FILE_PATH}
    RETURN  ${body.get('metadata').get('name')}

Get Test App Service Pod
    [Arguments]  ${pods}  ${app_name}
    ${test_app_pods}=  Get Object In Namespace By Mask  ${pods}  ${app_name}
    RETURN  ${test_app_pods}

Check Pods Count Is
    [Arguments]  ${count}  ${app_name}
    ${pods}=  Get Pods  ${namespace}
    ${test_app_pods}=  Get Test App Service Pod  ${pods}  ${app_name}
    ${list_len}=  Get List Length  ${test_app_pods}
    Convert To Integer  ${count}
    Should Be Equal As Integers  ${list_len}  ${count}

Check Status Of Pods
    [Arguments]  ${list_pods}
    FOR  ${pod}  IN  @{list_pods}
       ${state}=  Run Keyword And Return Status  Should Be Equal As Strings  ${pod.status.phase}  Running
       Should Be True  ${state}
       ...  Error! Following pod ${pod.metadata.name} has Failed status! Please, recheck pod status
    END
    RETURN  ${state}

Check That VMauth Is Presented In CR
     ${custom_resource}=  Get Custom Resource  monitoring.qubership.org/v1alpha1  PlatformMonitoring  ${namespace}  platformmonitoring
     ${flag}=  Check CR Service Exists  ${custom_resource.get('spec')}  ${vmauth-in-cr}  victoriametrics
     Log to console    vmauth ${flag}
     RETURN  ${flag}

Preparation Prometheus Session
    ${prometheus_url}=  Determine Protocol  ${prometheus_host}
    Create Session  prometheussession  ${prometheus_url}  verify=False

Preparation Victoriametrics Sessions With Oauth
    ${auth}=  Get Creadentials From Secret
    ${vmauth_url}=  Determine Protocol  ${vmauth_host}  ${auth}
    Create Session  vmsinglessession  ${vmauth_url}  auth=${auth}
    Create Session  vmagentsession  ${vmauth_url}  auth=${auth}

Preparation Victoriametrics Sessions Without Oauth
    ${vmsingle_url}=  Determine Protocol  ${vmsingle_host}
    ${vmagent_url}=   Determine Protocol  ${vmagent_host}
    Create Session  vmsinglessession  ${vmsingle_url}
    Create Session  vmagentsession  ${vmagent_url}

Preparation Victoriametrics Sessions
    ${vmauth}=  Check That VMauth Is Presented In CR
    Set Suite Variable  ${vmauth}
    Run Keyword If  '${vmauth}' == 'True'  Preparation Victoriametrics Sessions With Oauth
    ...  ELSE  Preparation Victoriametrics Sessions Without Oauth

Preparation Operator Session
    Run Keyword If  '${OPERATOR}' == 'prometheus-operator'
    ...  Preparation Prometheus Session
    Run Keyword If  '${OPERATOR}' == 'victoriametrics-operator'
    ...  Preparation Victoriametrics Sessions
