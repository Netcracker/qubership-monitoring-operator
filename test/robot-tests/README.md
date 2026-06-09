# Monitoring tests

## Test cases

1. [Grafana](src/tests/grafana)
   * Verify Login To Grafana
   * Create object grafana dashboard
   * Check that dashboard is appear in Grafana CR and includes name, namespace, uid
   * Update GrafanaDashboard And Check It Updated
   * Delete Dashboard and check it deleted in Grafana
2. [Smoke tests](src/tests/smoke-test/)
   * Check pods statuses of components managed by operator
      * Check Grafana Deployment Pods Are Running
      * Check Grafana Operator Pods Are Running
      * Check Monitoring Operator Pods Are Running
      * Check Prometheus Operator Pods Are Running
      * Check Alertmanager Pods Are Running
      * Check Prometheus Pods Are Running
      * Check Node Exporter Pods Are Running
      * Check Kube State Metrics Pods Are Running
      * Check Pushgateway Pods Are Running
   * Check pods statuses of components non managed by operator
      * Check Configurations Streamer Deployment Pods Are Running
      * Check Version Exporter Deployment Pods Are Running
      * Check Graphite Remote Adapter Deployment Pods Are Running
      * Check Cert Exporter Deployment Pods Are Running
      * Check Cloudwatch Exporter Deployment Pods Are Running
      * Check Blackbox Exporter Deployment Pods Are Running
      * Check VMAgent Deployment Pods Are Running
      * Check Prometheus Adapter Deployment Pods Are Running
      * Check Prometheus Adapter Operator Deployment Pods Are Running
      * Check Promxy Deployment Pods Are Running
      * Check Promitor Agent Scraper Deployment Pods Are Running
      * Check Network Latency Exporter Pods Are Running
   * Check Prometheus status
      * Check Status Of Prometheus
      * Check Status Of Prometheus Api
   * Check statuses of Ingresses/Routes
      * Check Prometheus Route Status
      * Check Status Of Prometheus Web Api
      * Check AlertManager Route/Ingress Status
      * Check AlertManager UI Status
      * Check Grafana Route/Ingress Status
      * Check Grafana UI Status
      * Check Pushgateway Route/Ingress Status
      * Check Pushgateway UI Status
   * Check default metrics targets
      * Check api-server Prometheus Target metrics
      * Check Etcd Prometheus Target metrics
      * Check Kube Controller Manager Prometheus Target metrics
      * Check Kube Scheduler Prometheus Target metrics
      * Check Kubelet Prometheus Target metrics
      * Check Non Mandatory Prometheus Target metrics
   * Check Victoria metrics status
      * Check Victoriametrics Operator Pods Are Running
      * Check VMAgent Pods Are Running
      * Check VMSingle Pods Are Running
      * Check VMAlert Pods Are Running
      * Check VMAlertmanager Pods Are Running
   * Check statuses of Ingresses/Routes for Victoria metrics
      * Check VMAgent Route/Ingress Status
      * Check VMAgent UI Status
      * Check VMSingle Route/Ingress Status
      * Check VMSingle UI Status
      * Check VMAlert Route/Ingress Status
      * Check VMAlert UI Status
      * Check VMAlertmanager Route/Ingress Status
      * Check VMAlertmanager UI Status
   * Check default metrics targets for Victoria metrics
      * Check Etcd VMAgent Target Metrics
      * Check Kubelet VMAgent Target Metrics
      * Check Non Mandatory Victoriametrics Target Metrics
3. [Test application](src/tests/test-app)
   * Check Prometheus Target's Status, Url, Labels
   * Wait Until Metrics Of Test App Is Exist
   * Check Available Metrics Of Test App
   * Update Service Monitor Label To 'monitoring-test'
   * Target Of Test App Has Been Deleted And No Metrics Are Written
   * Return label: Update label in Custom Resource to 'monitoring'
4. [Autoscaling test application](src/tests/adapter)
   * Increase Autoscaler Replicas To 3 And Check Metric Is 5
   * Decrease Autoscaler Replicas To 2 And Check Metric Is 1
5. [Image tests](src/tests/image-tests)
   * Check Images in Deployments, StatefulSets and DaemonSets

### Tags

<!-- markdownlint-disable line-length -->
| Tag                         | Description                                                     | Tests                                                                                                                                                                                  |
| --------------------------- | --------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| full                        | Tag to run all tests                                            | all                                                                                                                                                                                    |
| smoke                       | Tag to run smoke tests                                          | [Smoke tests](src/tests/smoke-test/), [Image tests](src/tests/image-tests)                                                                     |
| grafana                     | Tag to run grafana tests                                        | [Grafana](src/tests/grafana)                                                                                                                                       |
| test-app                    | Tag to run tests of test application                            | [Test application](src/tests/test-app)                                                                                                                             |
| blackbox-exporter           | Tag to run blackbox-exporter tests, included in smoke           | [Smoke tests](src/tests/smoke-test/) - Check Blackbox Exporter Deployment Pods Are Running                                                                         |
| cert-exporter               | Tag to run cert-exporter tests, included in smoke               | [Smoke tests](src/tests/smoke-test/) - Check Cert Exporter Deployment Pods Are Running                                                                             |
| cloudwatch-exporter         | Tag to run cloudwatch-exporter tests, included in smoke         | [Smoke tests](src/tests/smoke-test/) - Check Cloudwatch Exporter Deployment Pods Are Running                                                                       |
| configurations-streamer     | Tag to run configurations-streamer tests, included in smoke     | [Smoke tests](src/tests/smoke-test/) - Check Configurations Streamer Deployment Pods Are Running                                                                   |
| graphite-remote-adapter     | Tag to run graphite-remote-adapter tests, included in smoke     | [Smoke tests](src/tests/smoke-test/) - Check Graphite Remote Adapter Deployment Pods Are Running                                                                   |
| network-latency-exporter    | Tag to run graphite-remote-adapter tests, included in smoke     | [Smoke tests](src/tests/smoke-test/) - Check Network Latency Exporter Pods Are Running                                                                             |
| prometheus-adapter          | Tag to run prometheus-adapter  tests, included in smoke         | [Smoke tests](src/tests/smoke-test/) - Check Prometheus Adapter Deployment Pods Are Running; [Autoscaling test application](src/tests/adapter) |
| prometheus-adapter-operator | Tag to run prometheus-adapter-operator tests, included in smoke | [Smoke tests](src/tests/smoke-test/) - Check Prometheus Adapter Operator Deployment Pods Are Running                                                               |
| promitor-agent-scraper      | Tag to run promitor-agent-scraper tests, included in smoke      | [Smoke tests](src/tests/smoke-test/) - Check Promitor Agent Scraper Deployment Pods Are Running                                                                    |
| promxy                      | Tag to run promxy tests, included in smoke                      | [Smoke tests](src/tests/smoke-test/) - Check Promxy Deployment Pods Are Running                                                                                    |
| version-exporter            | Tag to run version-exporter tests, included in smoke            | [Smoke tests](src/tests/smoke-test/) - Check Version Exporter Deployment Pods Are Running                                                                          |
| alertmanager                | Tag to run alertmanager tests, included in smoke                | [Smoke tests](src/tests/smoke-test/) - Check Alertmanager                                                                                                          |
| monitoring_images           | Tag to run hardcoded image tests, included in smoke             | [Image tests](src/tests/image-tests) - Check Images In Deployments, StatefulSets and DaemonSets                                                                    |
<!-- markdownlint-enable line-length -->

You can also exclude some tests from smoke/full using NOT before the tag. For example, `grafanaORsmokeNOTalertmanager`.
