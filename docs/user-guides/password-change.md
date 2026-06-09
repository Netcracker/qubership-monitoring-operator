This guide describes how to change passwords for Monitoring and all components in it during
installation and during the work.

# Set passwords during deployment

During deploy you can specify admin users and passwords for the next components:

* Grafana
* VMAuth

**Note:** Because the VMAuth used as a proxy to external access for all VictoriaMetrics components
it means that it's enough to change specific users and passwords only for VMAuth. Other components
have no auth inside the Cloud.

## Grafana deploy

In current versions the **single source of truth** for Grafana admin credentials is the
`grafana-admin-credentials` Secret. The Grafana pod reads admin user and password from
environment variables:

- `GF_SECURITY_ADMIN_USER`
- `GF_SECURITY_ADMIN_PASSWORD`

These variables are populated from the Secret.

During deployment you can configure these credentials via `values.yaml`:

```yaml
grafana:
  # Main way to specify admin credentials for Grafana.
  # These values are used to render the grafana-admin-credentials Secret when
  # grafana.disableDefaultAdminSecret=false (default).
  security:
    admin_user: admin
    admin_password: admin
```

Behaviour:

- By default `admin_user` and `admin_password` are set to `admin/admin`.
- If you set any of them to an empty string `""`, Monitoring Operator will:
    - on first installation: generate a random value and store it in the Secret;
    - on subsequent upgrades: keep the existing value from the Secret.
- For backward compatibility, if the old section
  `grafana.config.security.admin_user` / `grafana.config.security.admin_password`
  is specified, its non-empty values override `grafana.security.*` **only when
  rendering the Secret**. New configurations should use `grafana.security.*`.

The `grafana.disableDefaultAdminSecret` flag controls who is responsible for creating
the admin credentials Secret:

```yaml
grafana:
  # false (default): Helm renders grafana-admin-credentials from grafana.security.*
  #                  and passes it to Grafana via environment variables.
  # true:            user is fully responsible for creating the Secret
  #                  {grafana-name}-admin-credentials with the required keys:
  #                  GF_SECURITY_ADMIN_USER and GF_SECURITY_ADMIN_PASSWORD.
  disableDefaultAdminSecret: false
```

- When `disableDefaultAdminSecret=false` (default), Helm always creates/updates
  the `grafana-admin-credentials` Secret based on values from `grafana.security.*`
  (or legacy `config.security.*`, if present).
- When `disableDefaultAdminSecret=true`, Helm does **not** create the Secret.
  Grafana and Grafana Operator read credentials from the Secret only if it exists.
  If the Secret is absent, Grafana falls back to its built-in default (`admin/admin`),
  so login is still possible.

Other external users and their passwords can't be set during deploy. During deploy you can specify
only which auth provides will use in Grafana.

If Grafana was configured use a Basic Auth so you can use the official guide to change their
passwords
[https://grafana.com/docs/grafana/latest/administration/user-management/user-preferences/](https://grafana.com/docs/grafana/latest/administration/user-management/user-preferences/).

**Note:** Please pay attention that if you are using OAuth2 or LDAP, or other external identity providers
you need to manage users and their passwords in these identity providers.

## VMAuth deploy

To specify VMAuth user during deploy you have to to add following in the deployment parameters:

```yaml
victoriametrics:
  vmUser:
    install: true
    username: prometheus
    password: prometheus
```

Also, it's possible to specify the password of user from the Secret.

**Warning!** The Secret with a password must be pre-created before deploy.

```yaml
victoriametrics:
  vmUser:
    install: true
    username: prometehus
    passwordRef:
      name: vmauth-secret  # the Secret name
      key: pass            # the key name inside the Secret
```

# Change passwords after deploy

This section describes how to change user credentials in runtime.

**Note:** After you will change credentials please do not forget to change them in the CMDB parameters.

## Grafana admin password change

To change Grafana's admin password you need to edit `grafana-admin-credentials` secret.

Find it in the namespace with Monitoring, for example using a command:

```bash
kubectl get secret -n <monitoring-namespace> grafana-admin-credentials
```

```bash
kubectl get secret -n <monitoring-namespace> grafana-admin-credentials

NAME                        TYPE     DATA   AGE
grafana-admin-credentials   Opaque   2      4h10m
```

And next need to edit it and change password:

```bash
> kubectl edit secrets -n <monitoring-namespace> grafana-admin-credentials
```

This opens your default editor and allows you to update the base64 encoded Secret values in the data field,
such as in the following example:

```yaml
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: v1
kind: Secret
metadata:
  name: grafana-admin-credentials
  ...
data:
  GF_SECURITY_ADMIN_PASSWORD: YWRtaW4=
  GF_SECURITY_ADMIN_USER: YWRtaW4=
type: Opaque
```

Update Base64 encoded password and save the file. Close the editor to update the secret.
Following message confirms the secret was edited successfully.

```bash
> kubectl edit secrets -n <monitoring-namespace> grafana-admin-credentials
secret/grafana-admin-credentials edited
```

### Release/0.57 or less

**NOTE:** If you use monitoring `release/0.57` version or less Grafana credentials are stored into grafana CR and
platform monitoring CR.

To change Grafana's admin password you need to edit PlatformMonitoring CR. Find it in the namespace with Monitoring,
for example using a command:

```shell
kubectl get -n <monitoring_namespace> platformmonitorings.monitoring.netcracker.com
```

Usually it has a name platformmonitoring:

```shell
kubectl get -n monitoring platformmonitorings.monitoring.netcracker.com

NAME                 AGE
platformmonitoring   11d
```

And next need to edit it and change password:

```yaml
grafana:
  config:
    security:
      admin_user: admin
      admin_password: admin
```

Monitoring-operator will start reconcile process, update Grafana CR and re-create grafana pod with new credentials.

## VMAuth password change

To change VMAuth credentials in runtime you need to edit PlatformMonitoring CR or a secret with a password.

In case if password specified in the CR, find it in the namespace with Monitoring, for example using a command:

```bash
kubectl get -n <monitoring_namespace> platformmonitorings.monitoring.netcracker.com
```

usually it have a name `platformmonitoring`:

```bash
❯ kubectl get -n monitoring platformmonitorings.monitoring.netcracker.com
NAME                 AGE
platformmonitoring   11d
```

And next need to edit it and change password:

```yaml
victoriametrics:
  vmUser:
    install: true
    username: prometehus
    password: prometheus
```

In case if password specified in the Secret, you need to find this Secret and change content in it.

The name of the secret you can find in the CMDB or PlatformMonitoring CR:

```yaml
victoriametrics:
  vmUser:
    passwordRef:
      name: vmauth-secret  # the Secret name
      key: pass            # the key name inside the Secret
```

After it edit a Secret:

```bash
kubectl edit -n <monitoring_namespace> secret <secret_name>
```

**Note:** Please keep in mind, that all values in the Secret stored in `base64` encoding. And before edit or save
data you must encode them in `base64`. In Linux you can use a command:

```bash
echo -n "<password>" | base64
```
