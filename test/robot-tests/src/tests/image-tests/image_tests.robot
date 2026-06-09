*** Settings ***
Library  String
Library  Collections
Library  PlatformLibrary  managed_by_operator=true

*** Variables ***
${NAMESPACE}      %{NAMESPACE}

*** Keywords ***
Compare Images From Resources With Dd
    [Arguments]  ${dd_images}
    ${stripped_resources}=  Strip String  ${dd_images}  characters=,  mode=right
    @{list_resources} =  Split String	${stripped_resources} 	,
    FOR  ${resource}  IN  @{list_resources}
      ${type}  ${name}  ${container_name}  ${image}=  Split String	${resource}
      ${resource_image}=  Get Resource Image  ${type}  ${name}  ${NAMESPACE}  ${container_name}
      Should Be Equal  ${resource_image}  ${image}
    END

*** Test Cases ***
Test Hardcoded Images
    [Tags]  monitoring_images  smoke  full
    ${dd_images}=  Get Dd Images From Config Map  tests-config  ${NAMESPACE}
    Skip If  '${dd_images}' == '${None}'  There is no deployDescriptor, not possible to check case!
    Compare Images From Resources With Dd  ${dd_images}

