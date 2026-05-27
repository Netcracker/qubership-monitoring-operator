# Password change guide

This guide describes how to change passwords for Monitoring and all components in it during
installation and during the work.

## Set passwords during deployment

During deploy you can specify admin users and passwords for the next components:

* Grafana
* VMAuth

**Note:** Because the VMAuth used as a proxy to external access for all VictoriaMetrics components
it means that it's enough to change specific users and passwords only for VMAuth. Other components
have no auth inside the Cloud.

### Grafana deploy

In current versions the **single source of truth** for Grafana admin credentials is the
`{grafana-name}-admin-credentials` Secret (default name: `grafana-admin-credentials`).
The Grafana pod reads admin user and password from files mounted from this Secret via
Grafana file provider (`$__file{...}`) in `grafana.ini`:

- `/etc/grafana-admin/GF_SECURITY_ADMIN_USER`
- `/etc/grafana-admin/GF_SECURITY_ADMIN_PASSWORD`

Secret keys remain `GF_SECURITY_ADMIN_USER` and `GF_SECURITY_ADMIN_PASSWORD`.

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
- If you set any of them to an empty string `""`, Helm will:
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
  #                  and mounts it into the Grafana pod.
  # true:            user is fully responsible for creating the Secret
  #                  {grafana-name}-admin-credentials with the required keys:
  #                  GF_SECURITY_ADMIN_USER and GF_SECURITY_ADMIN_PASSWORD.
  disableDefaultAdminSecret: false
```

- When `disableDefaultAdminSecret=false` (default), Helm creates/updates
  the `grafana-admin-credentials` Secret based on values from `grafana.security.*`
  (or legacy `config.security.*`, if present). Monitoring Operator instructs
  grafana-operator not to auto-generate its own admin secret.
- When `disableDefaultAdminSecret=true`, Helm does **not** create the Secret.
  Grafana reads credentials from the Secret only if it exists.
  If the Secret is absent, Grafana falls back to its built-in default (`admin/admin`),
  so login is still possible.

#### First start vs secret change (operator behaviour)

| Phase                           | Who applies credentials                                                                  | Monitoring Operator action                                                                                                                                                           |
|---------------------------------|------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **First install**               | Grafana reads mounted secret files on startup and creates the admin user in its database | Does **not** run `grafana cli admin reset-admin-password` (Grafana CR does not exist yet on the first credential check)                                                              |
| **Secret changed after deploy** | Secret is the source of truth                                                            | Detects checksum change, updates Grafana CR pod-template annotation `checksum/admin-secret` (rolling restart), then runs `grafana cli admin reset-admin-password` in the Grafana pod |

On deployments with a **Persistent Volume**, Grafana ignores `admin_password` from config once the admin user
already exists in the database. The CLI reset step is required in that case. On deployments without PV,
the rolling restart alone is often sufficient; the CLI reset is still executed and is idempotent.

**Reconcile interval:** credential sync runs on the next PlatformMonitoring reconcile cycle (default ~60 seconds
after the Secret is updated).

Other external users and their passwords can't be set during deploy. During deploy you can specify
only which auth provides will use in Grafana.

If Grafana was configured use a Basic Auth so you can use the official guide to change their
passwords
[https://grafana.com/docs/grafana/latest/administration/user-management/user-preferences/](https://grafana.com/docs/grafana/latest/administration/user-management/user-preferences/).

**Note:** Please pay attention that if you are using OAuth2 or LDAP, or other external identity providers
you need to manage users and their passwords in these identity providers.

### VMAuth deploy

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

## Change passwords after deploy

This section describes how to change user credentials in runtime.

**Note:** After you will change credentials please do not forget to change them in the CMDB parameters.

### Grafana admin password change

To change Grafana's admin password in runtime, edit the `grafana-admin-credentials` Secret.
This is the **only supported** way to keep Kubernetes and the running Grafana instance in sync.

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

Update base64 encoded password and save the file. Close the editor to update the secret.
Following message confirms the secret was edited successfully.

```bash
> kubectl edit secrets -n <monitoring-namespace> grafana-admin-credentials
secret/grafana-admin-credentials edited
```

#### What happens after the Secret is updated

Monitoring Operator reconciles PlatformMonitoring on a timer (default every 60 seconds). When it
detects that the Secret checksum differs from `checksum/admin-secret` on the Grafana CR pod template:

1. Updates the Grafana CR with the new checksum annotation (grafana-operator performs a rolling restart
   of the Grafana Deployment).
2. Waits until Grafana pods are ready.
3. Executes `grafana cli admin reset-admin-password <password-from-secret>` inside a running Grafana pod.

Expected log messages (`monitoring-operator` deployment logs):

```text
Admin credentials secret changed; Grafana credentials will be reset
Waiting for Grafana pods readiness before credential reset
Grafana admin credentials reset successfully
```

Verify the new password (replace namespace and password as needed):

```bash
kubectl port-forward -n <monitoring-namespace> svc/grafana-service 3000:3000
curl -u admin:<new-password> http://localhost:3000/api/org
```

You can also confirm that the checksum annotation on the Grafana CR was updated:

```bash
kubectl get grafana grafana -n <monitoring-namespace> \
  -o jsonpath='{.spec.deployment.spec.template.metadata.annotations.checksum\/admin-secret}{"\n"}'
```

If the Secret is **not** changed, the operator does not run credential reset on periodic reconciles.

#### Release/0.57 or less

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

### VMAuth password change

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
