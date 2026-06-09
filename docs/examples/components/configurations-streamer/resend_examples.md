This section provides details and examples of `configurations-streamer` resend APIs.

## HTTP Basic Authentication

All APIs will be authorized by the application using HTTP Basic Authentication.

Example:

```yaml
curl ...... -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

## API

* POST /api/v1/send

Example:

```yaml
curl -X POST http://localhost:8282/api/v1/send -d '{"entityTypes" : ["grafanadashboards","prometheusrules"],"namespaces" : ["monitoring","logging"],"targets" : {"grafana" : {"hosts" : ["https://10.101.17.197/grafana"]},"alertManager" : {"hosts" : ["http://10.101.17.197:6575"]}}}' -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, dashboards and Rules are resent only from monitoring and logging namespaces on specified targets.

```yaml
curl -X POST http://localhost:8282/api/v1/send -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, dashboards and Rules are resent on all namespaces on all enpoints.

* POST /api/v1/send/prometheusrules

Example:

```yaml
curl -X POST http://localhost:8282/api/v1/send/prometheusrules -d '{"entityTypes" : ["grafanadashboards","prometheusrules"],"namespaces" : ["monitoring","logging"],"targets" : {"grafana" : {"hosts" : ["https://10.101.17.197/grafana"]},"alertManager" : {"hosts" : ["http://10.101.17.197:6575"]}}}' -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, only rules will be resent from monitoring and logging namespaces on specified targets.

* POST /api/v1/send/grafanadashboards

Example:

```yaml
curl -X POST http://localhost:8282/api/v1/send/grafanadashboards -d '{"entityTypes" : ["grafanadashboards","prometheusrules"],"namespaces" : ["monitoring","logging"],"targets" : {"grafana" : {"hosts" : ["https://10.101.17.197/grafana"]},"alertManager" : {"hosts" : ["http://10.101.17.197:6575"]}}}' -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, only grafana dashboards will be resent from monitoring and logging namespaces on specified targets.

* GET /api/v1/send/grafanadashboards/api?endpoint=https%3A%2F%2F10.101.17.197%2Fgrafana

Note that the `endpoint` must be escaped. Multiple endpoints may be specified using multiple query parameters.

Example:

```yaml
curl -X GET http://localhost:8282/api/v1/send/grafanadashboards/api?endpoint=https%3A%2F%2F10.101.17.197%2Fgrafana -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, only dashboards are resent on all namespaces on API endpoint.

* GET /api/v1/send/grafanadashboards/ftp?endpoint=10.236.166.102%3A21

Note that the `endpoint` must be escaped. Multiple endpoints may be specified using multiple query parameters.

Example:

```yaml
curl -X GET http://localhost:8282/api/v1/send/grafanadashboards/ftp?endpoint=10.236.166.102%3A21 -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, only dashboards are resent on all namespaces on FTP endpoint 10.236.166.102:21.

* GET /api/v1/send/prometheusrules/rest?endpoint=http%3A%2F%2F10.101.17.197%3A6575

Note that the `endpoint` must be escaped. Multiple endpoints may be specified using multiple query parameters.

Example:

```yaml
curl -X GET http://localhost:8282/api/v1/send/prometheusrules/rest?endpoint=http%3A%2F%2F10.101.17.197%3A6575 -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, only rules are resent on all namespaces on REST endpoint.

* GET /api/v1/send/prometheusrules/ftp?endpoint=10.236.166.102%3A21

Note that the `endpoint` must be escaped. Multiple endpoints may be specified using multiple query parameters.

Example:

```yaml
curl -X GET http://localhost:8282/api/v1/send/prometheusrules/ftp?endpoint=10.236.166.102%3A21 -H 'Authorization: Basic YWRtaW46YWRtaW4='
```

In this example, only rules are resent on all namespaces on FTP endpoint 10.236.166.102:21.
