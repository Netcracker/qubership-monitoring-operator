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

To specify Grafana user and password in the deployment parameters you need to add following
in the deployment parameters:

```yaml
grafana:
  security:
    admin_user: admin
    admin_password: admin
```

This parameter allows configuring the Grafana admin user that grafana-operator uses to provision
Grafana, dashboards, data sources and other settings. Grafana operator creates `grafana-admin-credentials` by default,
but with a random password. We override these values during deployment using above configuration.

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
kubectl get -n <monitoring_namespace> platformmonitorings.monitoring.qubership.org
```

Usually it has a name platformmonitoring:

```shell
kubectl get -n monitoring platformmonitorings.monitoring.qubership.org

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
kubectl get -n <monitoring_namespace> platformmonitorings.monitoring.qubership.org
```

usually it have a name `platformmonitoring`:

```bash
❯ kubectl get -n monitoring platformmonitorings.monitoring.qubership.org
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
