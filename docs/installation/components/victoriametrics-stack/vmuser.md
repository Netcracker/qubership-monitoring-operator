#### vmuser

<!-- markdownlint-disable line-length -->
| Field            | Description                                                                                  | Scheme                | Required |
| ---------------- | -------------------------------------------------------------------------------------------- | --------------------- | -------- |
| install          | Allows to enable or disable deploy vmuser.                                                   | boolean               | false    |
| image            | A Docker image to deploy the vmuser.                                                         | string                | false    |
| paused           | Set paused to reconciliation for vmuser.                                                     | boolean               | false    |
| username         | UserName basic auth username for accessing protected endpoint. Default value: `admin`        | *string               | false    |
| password         | Password basic auth password for accessing protected endpoint. Default value: `admin`        | *string               | false    |
| passwordRef      | PasswordRef allows fetching password from user-create secret by its name and key.            | *v1.SecretKeySelector | false    |
| tokenRef         | TokenRef allows fetching token from user-created secrets by its name and key.                | *v1.SecretKeySelector | false    |
| generatePassword | GeneratePassword instructs operator to generate password for user if spec.password if empty. | *v1.SecretKeySelector | false    |
| bearerToken      | BearerToken Authorization header value for accessing protected endpoint.                     | *string               | false    |
| targetRefs       | TargetRefs - reference to endpoints, which user may access.                                  | []v1beta1.TargetRef   | false    |
<!-- markdownlint-enable line-length -->

```yaml
  vmUser:
    install: true
    paused: false
    username: vmauth
    passwordRef:
      key: pass
      name: vmauth-secret
```

If targetRefs is empty - targetRefs for vmalert, vmalertmanager, vmagent and vmsingle will be automatically created.
Default paths for vmalert:

* /api/v1/rules
* /api/v1/alerts
* /api/v1/alert.*
* /vmalert.*

Default paths for vmalertmanager:

* /api/v2/alerts.*
* /api/v2/receivers.*
* /api/v2/silences.*
* /api/v2/status.*

**TODO:** At this moment, you can navigate to the VMAlertManager UI through its Ingress, without using VMAuth, because
VMAlertManager is unavailable via VMAuth Ingress. We need to configure separated routing for it based on another Ingress
or on changing the UI with references to VM components.

Default paths for vmagent:

* /target.*
* /service-discovery.*
* /api/v1/write
* /api/v1/import.*
* /api/v1/target.*

Default paths for vmsingle:

* /vmui.*
* /graph.*
* /api/v1/label.*
* /api/v1/query.*
* /api/v1/rules
* /api/v1/alerts
* /api/v1/metadata.*
* /api/v1/format.*
* /api/v1/series.*
* /api/v1/status.*
* /api/v1/export.*
* /api/v1/admin/tsdb.*
* /prometheus/graph.*
* /prometheus/api/v1/label.*
* /graphite.*
* /prometheus/api/v1/query.*
* /prometheus/api/v1/rules
* /prometheus/api/v1/alerts
* /prometheus/api/v1/metadata
* /prometheus/api/v1/series.*
* /prometheus/api/v1/status.*
* /prometheus/api/v1/export.*
* /prometheus/federate
* /prometheus/api/v1/admin/tsdb.*

In case of OAuth + vmauth `Basic Auth` is supported only.
If OAuth is going to be installed - username (field "username: vmauth" in the example above) is to match the OAuth one,
auth `basicAuthPwd` field is to contain vmauth password.
Here is the config example.

> **Note:**
>
> If you skip parameters `username` and `password` they should be automatically set with default values `admin/admin`.

```yaml
   vmAuth:
     install: true
     ingress:
       install: true
       host: vmauth.test.org
     paused: false
     secretName: vmauth-secret
     extraVarsSecret:
       pass: vmauth
     securityContext:
       runAsUser: 2000
       fsGroup: 2000
   vmUser:
     install: true
     paused: false
     username: test
     passwordRef:
       key: pass
       name: vmauth-secret
 auth:
   clientId: vmtest
   clientSecret: vmtest
   loginUrl: https://idp.test.org/realms/myrealm/protocol/openid-connect/auth
   tokenUrl: https://idp.test.org/realms/myrealm/protocol/openid-connect/token
   userInfoUrl: https://idp.test.org/realms/myrealm/protocol/openid-connect/userinfo
   tlsConfig:
     insecureSkipVerify: true
   basicAuthPwd: vmauth
 oAuthProxy:
   image: oauth2-proxy/oauth2-proxy:v7.4.0
```

Default paths for vmselect:

* `/select/.*`


