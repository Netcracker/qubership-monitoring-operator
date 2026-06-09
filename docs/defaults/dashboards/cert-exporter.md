# SSL/TLS Certificates

Shows certificate expiration dates

## Tags

* `k8s`
* `certificates`

## Panels

### Certificates from files

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total Unique Certificates from Files | Total Unique Certificates from Files |  |  |
| Certificates from Files | Detailed information about certificates from files |  |  |
| Expiring Soon (Files) | Total count of certificates from files that expire within 30 days | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Expired (Files) | Total count of certificates from files that already expired | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Certificates from kubeconfig

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total Unique Certificates from Kubeconfig | Total Unique Certificates from Kubeconfig files |  |  |
| Certificates from Kubeconfig | Detailed information about certificates from kubeconfig files |  |  |
| Expiring Soon (Kubeconfig) | Total count of certificates from kubeconfig that expire within 30 days | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Expired (Kubeconfig) | Total count of certificates from kubeconfig that already expired | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Certificates from secrets

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total Unique Certificates from Secrets | Total Unique Certificates from Secrets |  |  |
| Certificates from Secrets | Detailed information about certificates from secrets |  |  |
| Expiring Soon (Secrets) | Total count of certificates from secrets that expire within 30 days | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Expired (Secrets) | Total count of certificates from secrets that already expired | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Exporter errors

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Errors[1h] | Errors that occurred during the exporter's work in the last 1 hour. See pod logs for details. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
