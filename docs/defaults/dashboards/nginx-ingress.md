# NGINX Ingress controller

Ingress-nginx supports a rich collection of prometheus metrics. If you have prometheus and grafana installed on your
cluster then prometheus will already be scraping this data due to the scrape annotation on the deployment.

## Tags

* `nginx`
* `k8s`
* `ingress`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Controller Request Volume | Show number of requests per second | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Controller Connections | Show average number of client connection | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Controller Success Rate (non-4&#124;5xx responses) | Show percentage successfully requests | Default:<br/>Mode: percentage<br/>Level 1: 90<br/>Level 2: 95<br/>Level 3: 99<br/><br/> |  |
| Config Reloads | Show number of Ingress controller reload operations | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Config Failed | Shows that last configuration failed | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Ingress Request Volume | Show the total number of client requests by ingress |  |  |
| Ingress Success Rate (non-4&#124;5xx responses) | Show percentage successfully requests by ingress |  |  |
| Network I/O pressure | Show request size (including request line, header, and request body) |  |  |
| Memory Usage | Show the number of bytes of memory in use  |  |  |
| CPU Usage | Show cpu usage  |  |  |
| Ingress Percentile Response Times and Transfer Rates | Show the requests processing time in milliseconds (le = 0.5, 09, 0.99) an requests size by ingress | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Ingress Certificate Expiry | Shows time until the SSL Certificate expires |  |  |
<!-- markdownlint-enable line-length -->
