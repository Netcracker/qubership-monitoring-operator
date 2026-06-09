# Blackbox Probes

Show state of probes of any external or internal URLs via Blackbox exporter. It check connectivity, response code and
time, and latency. Also it check certificates on URLs and expiration date of certificates.

## Tags

* `probe`
* `k8s`
* `blackbox`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Global Probe Duration | Shows the probe duration by each target | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### All targets SSL expiry

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| SSL Expiry | Shows earliest SSL cert expiry date for all targets. Please pay attention that SSL expiry shows only for targets with SSL encryption.  | Default:<br/>Mode: absolute<br/>Level 1: 1209600<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Blackbox self metrics

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Last successful config reload | Show last successful config reload |  |  |
| Version | Show Blackbox exporter version |  |  |
| Go Version | Show Golang version which use to compile Blackbox exporter |  |  |
| Unknown modules | Show count of unknown modules. This metric increase in case when Probe or Job try to execute probe with using unknown module name.<br/>`0` is a good value. |  |  |
| Config loaded state | Show state of config reload. `1` means that config was successfully reloaded. `0` means that config contains errors and wasn't reload. |  |  |
<!-- markdownlint-enable line-length -->

### $target status

**Row $target status is multiplied by parameter `target`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Status | Shows the status of probe | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| HTTP Duration | Shows the duration of HTTP request by phase | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Probe Duration | Shows how long the probe took to complete in seconds | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| HTTP Status Code | Shows response HTTP status code | Default:<br/>Mode: absolute<br/>Level 1: 201<br/>Level 2: 399<br/><br/> |  |
| HTTP Version | Shows the version of HTTP of the probe response | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| SSL Expiry | Shows earliest SSL cert expiry date | Default:<br/>Mode: absolute<br/>Level 1: 1209600<br/><br/> |  |
| Average Probe Duration | Shows average duration of the probe | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Average DNS Lookup | Shows average time taken for probe DNS lookup | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| SSL | Indicates if SSL was used for the final redirect | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->
