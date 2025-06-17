This document describes the metrics list and how to collect them from Etcd.

# Metrics

Etcd already exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint | Need Exporter? | Auth?                | Is Exporter Third Party? |
| ---------- | ------------ | ---------------- | -------------- | -------------------- | ------------------------ |
| Prometheus | `2379`       | `/metrics`       | No             | Require, certificate | N/A                      |

## How to Collect

Metrics expose on port `2379` and endpoint `/metrics`. By default, Etcd has authentication by certificates.

By default Etcd has no Service, so to collect metrics from it need to create a Service in the namespace where
Etcd deployed (in Kubernetes usually it's `kube-system`). Config of Service:

```yaml
kind: Service
apiVersion: v1
metadata:
  name: etcd
  labels:
    k8s-app: etcd
spec:
  ports:
    - name: metrics
      protocol: TCP
      port: 2379
      targetPort: 2379
  selector:
    component: etcd
  clusterIP: None
  clusterIPs:
    - None
  type: ClusterIP
  sessionAffinity: None
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  internalTrafficPolicy: Cluster
```

Config `ServiceMonitor` for `prometheus-operator` to collect Etcd metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
  name: monitoring-etcd-service-monitor
spec:
  endpoints:
    - interval: 30s
      port: metrics
      scheme: https
      scrapeTimeout: 10s
      tlsConfig:
        caFile: /etc/prometheus/secrets/kube-etcd-client-certs/etcd-client-ca.crt
        certFile: /etc/prometheus/secrets/kube-etcd-client-certs/etcd-client.crt
        keyFile: /etc/prometheus/secrets/kube-etcd-client-certs/etcd-client.key
  jobLabel: k8s-app
  namespaceSelector:
    matchNames:
      - kube-system
  selector:
    matchLabels:
      k8s-app: etcd
```

### How to find certificates for Etcd

If you configure Etcd monitoring manually, so you need to fill empty Secret for Monitoring with found Etcd
certificates. For example:

```yaml
apiVersion: v1
kind: Secret
metadata:
  labels:
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: prometheus-operator
    app.kubernetes.io/name: kube-etcd-client-certs
  name: kube-etcd-client-certs
  namespace: prometheus-operator
type: Opaque
data:
  etcd-client-ca.crt: <your_etcd-client-ca.crt>
  etcd-client.crt: <your_etcd-client.crt>
  etcd-client.key: <your_etcd-client.key>
```

To find certificate in different Kubernetes and OpenShift version exist different ways.

#### For Kubernetes 1.19-?

```bash
# etcd-client-ca.crt
kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | kubectl exec $(awk '/etcd/{print $1; exit}') -n kube-system -- echo $(</etc/kubernetes/pki/etcd/ca.crt) | base64 | sed ':a; /$/N; s/\n//; ta'

# etcd-client.crt
kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | kubectl exec $(awk '/etcd/{print $1; exit}') -n kube-system --  echo $(</etc/kubernetes/pki/etcd/peer.crt) | base64 | sed ':a; /$/N; s/\n//; ta'

# etcd-client.key
kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | kubectl exec $(awk '/etcd/{print $1; exit}') -n kube-system --  echo $(</etc/kubernetes/pki/etcd/peer.key) | base64 | sed ':a; /$/N; s/\n//; ta'
```

**If the previous method didn't work:**

In some cases, the etcd container may not have a shell that would allow commands to be used directly through
`kubectl exec` command. In this case, the following way is recommended:

1. Open 2 linux terminal (the first will be used to access the terminal inside the etcd container,
the second to process the contents of the certificates).
2. Use commands below inside terminal №1 to open the shell into etcd container:

   ```bash
   etcd_pod=$(kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | awk '/etcd/{print $1; exit}')
   kubectl exec -it $etcd_pod -n kube-system -- sh
   ```

3. Inside the etcd container (terminal №1) run the following command to get the contents of the ca.crt certificate:

   ```bash
   echo $(</etc/kubernetes/pki/etcd/ca.crt)
   ```

4. Run the following command inside terminal №2. Replace `<certificate>` with output of previous step:

   ```bash
   echo "<certificate>" | sed 's/ /\n/g;s/-----BEGIN\nCERTIFICATE-----/-----BEGIN CERTIFICATE-----/;s/-----END\nCERTIFICATE-----/-----END CERTIFICATE-----/' | base64
   ```

5. You can put output of previous command in the `kube-etcd-client-certs` Secret as `etcd-client-ca.crt` value.
6. Get certificate from `peer.crt` file by running follow command inside terminal №1:

   ```bash
   echo $(</etc/kubernetes/pki/etcd/peer.crt)
   ```

7. Repeat step 4. You can put output of this command as `etcd-client.crt` in the Secret.
8. Get RSA private key from `peer.key` file by running follow command inside terminal №1:

   ```bash
   echo $(</etc/kubernetes/pki/etcd/peer.key)
   ```

9. Run the following command inside terminal №2. Replace `<certificate>` with output of previous step:

   ```bash
   echo "<certificate>" | sed 's/ /\n/g;s/-----BEGIN\nRSA\nPRIVATE\nKEY-----/-----BEGIN RSA PRIVATE KEY-----/;s/-----END\nRSA\nPRIVATE\nKEY-----/-----END RSA PRIVATE KEY-----/' | base64
   ```

10. You can put output of previous command as `etcd-client.key` value in the Secret.

#### For Kubernetes 1.15-1.18

```bash
# etcd-client-ca.crt
kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | kubectl exec $(awk '/etcd/{print $1; exit}') -n kube-system -- cat /etc/kubernetes/pki/etcd/ca.crt | base64 | sed ':a; /$/N; s/\n//; ta'

# etcd-client.crt
kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | kubectl exec $(awk '/etcd/{print $1; exit}') -n kube-system -- cat /etc/kubernetes/pki/etcd/peer.crt | base64 | sed ':a; /$/N; s/\n//; ta'

# etcd-client.key
kubectl get pods -n kube-system --no-headers --field-selector status.phase=Running | kubectl exec $(awk '/etcd/{print $1; exit}') -n kube-system -- cat /etc/kubernetes/pki/etcd/peer.key | base64 | sed ':a; /$/N; s/\n//; ta'
```

#### For OpenShift 4.x

```bash
# etcd-client-ca.crt
oc get configmap etcd-metric-serving-ca --namespace=openshift-etcd-operator -o json | jq -r '.data."ca-bundle.crt"' | base64

# etcd-client.crt
oc get secret etcd-metric-client --namespace=openshift-etcd-operator -o json | jq -r '.data."tls.crt"'

# etcd-client.key
oc get secret etcd-metric-client --namespace=openshift-etcd-operator -o json | jq -r '.data."tls.key"'
```

In additional for OpenShift 4.5 - 4.7 need to change settings into ServiceMonitor for etcd. It can be found in namespace
in which deployed monitoring and should have name like `<namesapce>-etcd-service-monitor`.

In this file need to change:

1. Change namespace, because in OpenShift 4.5-4.7 etcd was moved from `kube-system` to `openshift-etcd`:

    ```yaml
    spec:
      namespaceSelector:
        matchNames:
        - openshift-etcd
    ```

2. Change metrics port name, because in OpenShift 4.5-4.7 in service for etcd port was renamed to `etcd-metrics`:

    ```yaml
    spec:
      endpoints:
      - interval: 30s
         port: etcd-metrics
    ```

3. Check that both etcd pods and etcd service has label `k8s-app: etcd`. If this label not exists on one of both
   resources need eaither to add it. Or find other label which present on both resources and specify it
   in `matchSelector`:

   ```yaml
   spec:
      selector:
        matchLabels:
          k8s-app: etcd
   ```

#### For OpenShift 3.11

```bash
# etcd-client-ca.crt
oc get pods -n kube-system --no-headers --selector openshift.io/component=etcd | oc exec $(awk '/etcd/{print $1; exit}') -n kube-system -- cat /etc/etcd/ca.crt | base64 | sed ':a; /$/N; s/\n//; ta'

# etcd-client.crt
oc get pods -n kube-system --no-headers --selector openshift.io/component=etcd | oc exec $(awk '/etcd/{print $1; exit}') -n kube-system -- cat /etc/etcd/peer.crt | base64 | sed ':a; /$/N; s/\n//; ta'

# etcd-client.key
oc get pods -n kube-system --no-headers --selector openshift.io/component=etcd | oc exec $(awk '/etcd/{print $1; exit}') -n kube-system -- cat /etc/etcd/peer.key | base64 | sed ':a; /$/N; s/\n//; ta'
```

## Metrics List

```prometheus
# HELP etcd_cluster_version Which version is running. 1 for 'cluster_version' label with current cluster version
# TYPE etcd_cluster_version gauge
etcd_cluster_version{cluster_version="3.4"} 1
# HELP etcd_debugging_disk_backend_commit_rebalance_duration_seconds The latency distributions of commit.rebalance called by bboltdb backend.
# TYPE etcd_debugging_disk_backend_commit_rebalance_duration_seconds histogram
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.001"} 2.058614e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.002"} 2.058889e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.004"} 2.05903e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.008"} 2.059047e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.016"} 2.05905e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.032"} 2.05905e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.064"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.128"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.256"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="0.512"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="1.024"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="2.048"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="4.096"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="8.192"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_bucket{le="+Inf"} 2.059051e+06
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_sum 3.8554109809999733
etcd_debugging_disk_backend_commit_rebalance_duration_seconds_count 2.059051e+06
# HELP etcd_debugging_disk_backend_commit_spill_duration_seconds The latency distributions of commit.spill called by bboltdb backend.
# TYPE etcd_debugging_disk_backend_commit_spill_duration_seconds histogram
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.001"} 2.055145e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.002"} 2.057487e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.004"} 2.058237e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.008"} 2.058613e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.016"} 2.058844e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.032"} 2.058978e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.064"} 2.059039e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.128"} 2.059049e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.256"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="0.512"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="1.024"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="2.048"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="4.096"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="8.192"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_bucket{le="+Inf"} 2.059051e+06
etcd_debugging_disk_backend_commit_spill_duration_seconds_sum 205.06119727599784
etcd_debugging_disk_backend_commit_spill_duration_seconds_count 2.059051e+06
# HELP etcd_debugging_disk_backend_commit_write_duration_seconds The latency distributions of commit.write called by bboltdb backend.
# TYPE etcd_debugging_disk_backend_commit_write_duration_seconds histogram
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.001"} 0
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.002"} 0
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.004"} 487
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.008"} 1.282599e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.016"} 1.830559e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.032"} 1.961272e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.064"} 1.997001e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.128"} 2.021928e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.256"} 2.053133e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="0.512"} 2.057981e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="1.024"} 2.058874e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="2.048"} 2.059009e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="4.096"} 2.059047e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="8.192"} 2.059051e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_bucket{le="+Inf"} 2.059051e+06
etcd_debugging_disk_backend_commit_write_duration_seconds_sum 28946.043956317633
etcd_debugging_disk_backend_commit_write_duration_seconds_count 2.059051e+06
# HELP etcd_debugging_lease_granted_total The total number of granted leases.
# TYPE etcd_debugging_lease_granted_total counter
etcd_debugging_lease_granted_total 193481
# HELP etcd_debugging_lease_renewed_total The number of renewed leases seen by the leader.
# TYPE etcd_debugging_lease_renewed_total counter
etcd_debugging_lease_renewed_total 0
# HELP etcd_debugging_lease_revoked_total The total number of revoked leases.
# TYPE etcd_debugging_lease_revoked_total counter
etcd_debugging_lease_revoked_total 193620
# HELP etcd_debugging_lease_ttl_total Bucketed histogram of lease TTLs.
# TYPE etcd_debugging_lease_ttl_total histogram
etcd_debugging_lease_ttl_total_bucket{le="1"} 0
etcd_debugging_lease_ttl_total_bucket{le="2"} 0
etcd_debugging_lease_ttl_total_bucket{le="4"} 0
etcd_debugging_lease_ttl_total_bucket{le="8"} 0
etcd_debugging_lease_ttl_total_bucket{le="16"} 179589
etcd_debugging_lease_ttl_total_bucket{le="32"} 179589
etcd_debugging_lease_ttl_total_bucket{le="64"} 179589
etcd_debugging_lease_ttl_total_bucket{le="128"} 179589
etcd_debugging_lease_ttl_total_bucket{le="256"} 179589
etcd_debugging_lease_ttl_total_bucket{le="512"} 179589
etcd_debugging_lease_ttl_total_bucket{le="1024"} 179589
etcd_debugging_lease_ttl_total_bucket{le="2048"} 179589
etcd_debugging_lease_ttl_total_bucket{le="4096"} 193481
etcd_debugging_lease_ttl_total_bucket{le="8192"} 193481
etcd_debugging_lease_ttl_total_bucket{le="16384"} 193481
etcd_debugging_lease_ttl_total_bucket{le="32768"} 193481
etcd_debugging_lease_ttl_total_bucket{le="65536"} 193481
etcd_debugging_lease_ttl_total_bucket{le="131072"} 193481
etcd_debugging_lease_ttl_total_bucket{le="262144"} 193481
etcd_debugging_lease_ttl_total_bucket{le="524288"} 193481
etcd_debugging_lease_ttl_total_bucket{le="1.048576e+06"} 193481
etcd_debugging_lease_ttl_total_bucket{le="2.097152e+06"} 193481
etcd_debugging_lease_ttl_total_bucket{le="4.194304e+06"} 193481
etcd_debugging_lease_ttl_total_bucket{le="8.388608e+06"} 193481
etcd_debugging_lease_ttl_total_bucket{le="+Inf"} 193481
etcd_debugging_lease_ttl_total_sum 5.3538555e+07
etcd_debugging_lease_ttl_total_count 193481
# HELP etcd_debugging_mvcc_compact_revision The revision of the last compaction in store.
# TYPE etcd_debugging_mvcc_compact_revision gauge
etcd_debugging_mvcc_compact_revision 4.463841e+07
# HELP etcd_debugging_mvcc_current_revision The current revision of store.
# TYPE etcd_debugging_mvcc_current_revision gauge
etcd_debugging_mvcc_current_revision 4.4642331e+07
# HELP etcd_debugging_mvcc_db_compaction_keys_total Total number of db keys compacted.
# TYPE etcd_debugging_mvcc_db_compaction_keys_total counter
etcd_debugging_mvcc_db_compaction_keys_total 0
# HELP etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds Bucketed histogram of db compaction pause duration.
# TYPE etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds histogram
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="1"} 2901
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="2"} 3547
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="4"} 3662
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="8"} 3877
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="16"} 5434
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="32"} 6236
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="64"} 6319
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="128"} 6344
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="256"} 6399
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="512"} 6409
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="1024"} 6411
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="2048"} 6412
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="4096"} 6412
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_bucket{le="+Inf"} 6412
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_sum 65286
etcd_debugging_mvcc_db_compaction_pause_duration_milliseconds_count 6412
# HELP etcd_debugging_mvcc_db_compaction_total_duration_milliseconds Bucketed histogram of db compaction total duration.
# TYPE etcd_debugging_mvcc_db_compaction_total_duration_milliseconds histogram
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="100"} 1832
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="200"} 1924
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="400"} 1980
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="800"} 1990
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="1600"} 1991
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="3200"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="6400"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="12800"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="25600"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="51200"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="102400"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="204800"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="409600"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="819200"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_bucket{le="+Inf"} 1992
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_sum 139068
etcd_debugging_mvcc_db_compaction_total_duration_milliseconds_count 1992
# HELP etcd_debugging_mvcc_db_total_size_in_bytes Total size of the underlying database physically allocated in bytes.
# TYPE etcd_debugging_mvcc_db_total_size_in_bytes gauge
etcd_debugging_mvcc_db_total_size_in_bytes 2.26074624e+08
# HELP etcd_debugging_mvcc_delete_total Total number of deletes seen by this member.
# TYPE etcd_debugging_mvcc_delete_total counter
etcd_debugging_mvcc_delete_total 25008
# HELP etcd_debugging_mvcc_events_total Total number of events sent by this member.
# TYPE etcd_debugging_mvcc_events_total counter
etcd_debugging_mvcc_events_total 4.215829e+06
# HELP etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds Bucketed histogram of index compaction pause duration.
# TYPE etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds histogram
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="0.5"} 0
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="1"} 1222
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="2"} 1864
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="4"} 1951
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="8"} 1978
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="16"} 1984
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="32"} 1988
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="64"} 1990
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="128"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="256"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="512"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="1024"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="2048"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="4096"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_bucket{le="+Inf"} 1992
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_sum 3381
etcd_debugging_mvcc_index_compaction_pause_duration_milliseconds_count 1992
# HELP etcd_debugging_mvcc_keys_total Total number of keys.
# TYPE etcd_debugging_mvcc_keys_total gauge
etcd_debugging_mvcc_keys_total 2733
# HELP etcd_debugging_mvcc_pending_events_total Total number of pending events to be sent.
# TYPE etcd_debugging_mvcc_pending_events_total gauge
etcd_debugging_mvcc_pending_events_total 0
# HELP etcd_debugging_mvcc_put_total Total number of puts seen by this member.
# TYPE etcd_debugging_mvcc_put_total counter
etcd_debugging_mvcc_put_total 3.016864e+06
# HELP etcd_debugging_mvcc_range_total Total number of ranges seen by this member.
# TYPE etcd_debugging_mvcc_range_total counter
etcd_debugging_mvcc_range_total 1.2084894e+07
# HELP etcd_debugging_mvcc_slow_watcher_total Total number of unsynced slow watchers.
# TYPE etcd_debugging_mvcc_slow_watcher_total gauge
etcd_debugging_mvcc_slow_watcher_total 0
# HELP etcd_debugging_mvcc_txn_total Total number of txns seen by this member.
# TYPE etcd_debugging_mvcc_txn_total counter
etcd_debugging_mvcc_txn_total 2144
# HELP etcd_debugging_mvcc_watch_stream_total Total number of watch streams.
# TYPE etcd_debugging_mvcc_watch_stream_total gauge
etcd_debugging_mvcc_watch_stream_total 113
# HELP etcd_debugging_mvcc_watcher_total Total number of watchers.
# TYPE etcd_debugging_mvcc_watcher_total gauge
etcd_debugging_mvcc_watcher_total 113
# HELP etcd_debugging_server_lease_expired_total The total number of expired leases.
# TYPE etcd_debugging_server_lease_expired_total counter
etcd_debugging_server_lease_expired_total 48707
# HELP etcd_debugging_snap_save_marshalling_duration_seconds The marshalling cost distributions of save called by snapshot.
# TYPE etcd_debugging_snap_save_marshalling_duration_seconds histogram
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.001"} 377
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.002"} 378
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.004"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.008"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.016"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.032"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.064"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.128"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.256"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="0.512"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="1.024"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="2.048"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="4.096"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="8.192"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_bucket{le="+Inf"} 380
etcd_debugging_snap_save_marshalling_duration_seconds_sum 0.044155771
etcd_debugging_snap_save_marshalling_duration_seconds_count 380
# HELP etcd_debugging_snap_save_total_duration_seconds The total latency distributions of save called by snapshot.
# TYPE etcd_debugging_snap_save_total_duration_seconds histogram
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.001"} 0
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.002"} 0
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.004"} 0
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.008"} 181
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.016"} 343
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.032"} 368
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.064"} 373
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.128"} 376
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.256"} 380
etcd_debugging_snap_save_total_duration_seconds_bucket{le="0.512"} 380
etcd_debugging_snap_save_total_duration_seconds_bucket{le="1.024"} 380
etcd_debugging_snap_save_total_duration_seconds_bucket{le="2.048"} 380
etcd_debugging_snap_save_total_duration_seconds_bucket{le="4.096"} 380
etcd_debugging_snap_save_total_duration_seconds_bucket{le="8.192"} 380
etcd_debugging_snap_save_total_duration_seconds_bucket{le="+Inf"} 380
etcd_debugging_snap_save_total_duration_seconds_sum 4.608179085999998
etcd_debugging_snap_save_total_duration_seconds_count 380
# HELP etcd_debugging_store_expires_total Total number of expired keys.
# TYPE etcd_debugging_store_expires_total counter
etcd_debugging_store_expires_total 0
# HELP etcd_debugging_store_reads_total Total number of reads action by (get/getRecursive), local to this member.
# TYPE etcd_debugging_store_reads_total counter
etcd_debugging_store_reads_total{action="get"} 179667
etcd_debugging_store_reads_total{action="getRecursive"} 2
# HELP etcd_debugging_store_watch_requests_total Total number of incoming watch requests (new or reestablished).
# TYPE etcd_debugging_store_watch_requests_total counter
etcd_debugging_store_watch_requests_total 0
# HELP etcd_debugging_store_watchers Count of currently active watchers.
# TYPE etcd_debugging_store_watchers gauge
etcd_debugging_store_watchers 0
# HELP etcd_debugging_store_writes_total Total number of writes (e.g. set/compareAndDelete) seen by this member.
# TYPE etcd_debugging_store_writes_total counter
etcd_debugging_store_writes_total{action="set"} 3
# HELP etcd_disk_backend_commit_duration_seconds The latency distributions of commit called by backend.
# TYPE etcd_disk_backend_commit_duration_seconds histogram
etcd_disk_backend_commit_duration_seconds_bucket{le="0.001"} 0
etcd_disk_backend_commit_duration_seconds_bucket{le="0.002"} 0
etcd_disk_backend_commit_duration_seconds_bucket{le="0.004"} 247
etcd_disk_backend_commit_duration_seconds_bucket{le="0.008"} 1.25831e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="0.016"} 1.82696e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="0.032"} 1.960819e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="0.064"} 1.996831e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="0.128"} 2.021882e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="0.256"} 2.053122e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="0.512"} 2.05798e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="1.024"} 2.058874e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="2.048"} 2.059009e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="4.096"} 2.059047e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="8.192"} 2.059051e+06
etcd_disk_backend_commit_duration_seconds_bucket{le="+Inf"} 2.059051e+06
etcd_disk_backend_commit_duration_seconds_sum 29192.160583952755
etcd_disk_backend_commit_duration_seconds_count 2.059051e+06
# HELP etcd_disk_backend_defrag_duration_seconds The latency distribution of backend defragmentation.
# TYPE etcd_disk_backend_defrag_duration_seconds histogram
etcd_disk_backend_defrag_duration_seconds_bucket{le="0.1"} 0
etcd_disk_backend_defrag_duration_seconds_bucket{le="0.2"} 0
etcd_disk_backend_defrag_duration_seconds_bucket{le="0.4"} 3
etcd_disk_backend_defrag_duration_seconds_bucket{le="0.8"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="1.6"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="3.2"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="6.4"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="12.8"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="25.6"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="51.2"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="102.4"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="204.8"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="409.6"} 8
etcd_disk_backend_defrag_duration_seconds_bucket{le="+Inf"} 8
etcd_disk_backend_defrag_duration_seconds_sum 3.8479464640000005
etcd_disk_backend_defrag_duration_seconds_count 8
# HELP etcd_disk_backend_snapshot_duration_seconds The latency distribution of backend snapshots.
# TYPE etcd_disk_backend_snapshot_duration_seconds histogram
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.01"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.02"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.04"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.08"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.16"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.32"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="0.64"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="1.28"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="2.56"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="5.12"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="10.24"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="20.48"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="40.96"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="81.92"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="163.84"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="327.68"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="655.36"} 0
etcd_disk_backend_snapshot_duration_seconds_bucket{le="+Inf"} 0
etcd_disk_backend_snapshot_duration_seconds_sum 0
etcd_disk_backend_snapshot_duration_seconds_count 0
# HELP etcd_disk_wal_fsync_duration_seconds The latency distributions of fsync called by WAL.
# TYPE etcd_disk_wal_fsync_duration_seconds histogram
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.001"} 0
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.002"} 190419
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.004"} 2.267522e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.008"} 3.166919e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.016"} 3.574917e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.032"} 3.661494e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.064"} 3.688871e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.128"} 3.712177e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.256"} 3.740627e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.512"} 3.744871e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="1.024"} 3.745806e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="2.048"} 3.745957e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="4.096"} 3.745991e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="8.192"} 3.745994e+06
etcd_disk_wal_fsync_duration_seconds_bucket{le="+Inf"} 3.745994e+06
etcd_disk_wal_fsync_duration_seconds_sum 28709.93130914516
etcd_disk_wal_fsync_duration_seconds_count 3.745994e+06
# HELP etcd_grpc_proxy_cache_hits_total Total number of cache hits
# TYPE etcd_grpc_proxy_cache_hits_total gauge
etcd_grpc_proxy_cache_hits_total 0
# HELP etcd_grpc_proxy_cache_keys_total Total number of keys/ranges cached
# TYPE etcd_grpc_proxy_cache_keys_total gauge
etcd_grpc_proxy_cache_keys_total 0
# HELP etcd_grpc_proxy_cache_misses_total Total number of cache misses
# TYPE etcd_grpc_proxy_cache_misses_total gauge
etcd_grpc_proxy_cache_misses_total 0
# HELP etcd_grpc_proxy_events_coalescing_total Total number of events coalescing
# TYPE etcd_grpc_proxy_events_coalescing_total counter
etcd_grpc_proxy_events_coalescing_total 0
# HELP etcd_grpc_proxy_watchers_coalescing_total Total number of current watchers coalescing
# TYPE etcd_grpc_proxy_watchers_coalescing_total gauge
etcd_grpc_proxy_watchers_coalescing_total 0
# HELP etcd_mvcc_db_open_read_transactions The number of currently open read transactions
# TYPE etcd_mvcc_db_open_read_transactions gauge
etcd_mvcc_db_open_read_transactions 1
# HELP etcd_mvcc_db_total_size_in_bytes Total size of the underlying database physically allocated in bytes.
# TYPE etcd_mvcc_db_total_size_in_bytes gauge
etcd_mvcc_db_total_size_in_bytes 2.26074624e+08
# HELP etcd_mvcc_db_total_size_in_use_in_bytes Total size of the underlying database logically in use in bytes.
# TYPE etcd_mvcc_db_total_size_in_use_in_bytes gauge
etcd_mvcc_db_total_size_in_use_in_bytes 1.89222912e+08
# HELP etcd_mvcc_delete_total Total number of deletes seen by this member.
# TYPE etcd_mvcc_delete_total counter
etcd_mvcc_delete_total 25008
# HELP etcd_mvcc_hash_duration_seconds The latency distribution of storage hash operation.
# TYPE etcd_mvcc_hash_duration_seconds histogram
etcd_mvcc_hash_duration_seconds_bucket{le="0.01"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="0.02"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="0.04"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="0.08"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="0.16"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="0.32"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="0.64"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="1.28"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="2.56"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="5.12"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="10.24"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="20.48"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="40.96"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="81.92"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="163.84"} 0
etcd_mvcc_hash_duration_seconds_bucket{le="+Inf"} 0
etcd_mvcc_hash_duration_seconds_sum 0
etcd_mvcc_hash_duration_seconds_count 0
# HELP etcd_mvcc_hash_rev_duration_seconds The latency distribution of storage hash by revision operation.
# TYPE etcd_mvcc_hash_rev_duration_seconds histogram
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.01"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.02"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.04"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.08"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.16"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.32"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="0.64"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="1.28"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="2.56"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="5.12"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="10.24"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="20.48"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="40.96"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="81.92"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="163.84"} 0
etcd_mvcc_hash_rev_duration_seconds_bucket{le="+Inf"} 0
etcd_mvcc_hash_rev_duration_seconds_sum 0
etcd_mvcc_hash_rev_duration_seconds_count 0
# HELP etcd_mvcc_put_total Total number of puts seen by this member.
# TYPE etcd_mvcc_put_total counter
etcd_mvcc_put_total 3.016864e+06
# HELP etcd_mvcc_range_total Total number of ranges seen by this member.
# TYPE etcd_mvcc_range_total counter
etcd_mvcc_range_total 1.2084894e+07
# HELP etcd_mvcc_txn_total Total number of txns seen by this member.
# TYPE etcd_mvcc_txn_total counter
etcd_mvcc_txn_total 2144
# HELP etcd_network_active_peers The current number of active peer connections.
# TYPE etcd_network_active_peers gauge
etcd_network_active_peers{Local="b3e330dca330e585",Remote="b10c280072f91b6"} 1
etcd_network_active_peers{Local="b3e330dca330e585",Remote="ceb61f37051a244c"} 1
# HELP etcd_network_client_grpc_received_bytes_total The total number of bytes received from grpc clients.
# TYPE etcd_network_client_grpc_received_bytes_total counter
etcd_network_client_grpc_received_bytes_total 5.803630091e+09
# HELP etcd_network_client_grpc_sent_bytes_total The total number of bytes sent to grpc clients.
# TYPE etcd_network_client_grpc_sent_bytes_total counter
etcd_network_client_grpc_sent_bytes_total 5.8553765501e+10
# HELP etcd_network_disconnected_peers_total The total number of disconnected peers.
# TYPE etcd_network_disconnected_peers_total counter
etcd_network_disconnected_peers_total{Local="b3e330dca330e585",Remote="b10c280072f91b6"} 1
etcd_network_disconnected_peers_total{Local="b3e330dca330e585",Remote="ceb61f37051a244c"} 1
# HELP etcd_network_peer_received_bytes_total The total number of bytes received from peers.
# TYPE etcd_network_peer_received_bytes_total counter
etcd_network_peer_received_bytes_total{From="0"} 4.310496e+07
etcd_network_peer_received_bytes_total{From="b10c280072f91b6"} 5.188405015e+09
etcd_network_peer_received_bytes_total{From="ceb61f37051a244c"} 1.0219567106e+10
# HELP etcd_network_peer_round_trip_time_seconds Round-Trip-Time histogram between peers
# TYPE etcd_network_peer_round_trip_time_seconds histogram
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0001"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0002"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0004"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0008"} 1
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0016"} 6006
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0032"} 13582
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0064"} 17987
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0128"} 19685
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0256"} 38091
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.0512"} 39871
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.1024"} 39901
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.2048"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.4096"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="0.8192"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="1.6384"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="3.2768"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="b10c280072f91b6",le="+Inf"} 39912
etcd_network_peer_round_trip_time_seconds_sum{To="b10c280072f91b6"} 430.5429549429996
etcd_network_peer_round_trip_time_seconds_count{To="b10c280072f91b6"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0001"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0002"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0004"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0008"} 0
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0016"} 2401
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0032"} 11150
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0064"} 17731
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0128"} 19866
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0256"} 38809
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.0512"} 39891
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.1024"} 39908
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.2048"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.4096"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="0.8192"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="1.6384"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="3.2768"} 39912
etcd_network_peer_round_trip_time_seconds_bucket{To="ceb61f37051a244c",le="+Inf"} 39912
etcd_network_peer_round_trip_time_seconds_sum{To="ceb61f37051a244c"} 427.1415449249977
etcd_network_peer_round_trip_time_seconds_count{To="ceb61f37051a244c"} 39912
# HELP etcd_network_peer_sent_bytes_total The total number of bytes sent to peers.
# TYPE etcd_network_peer_sent_bytes_total counter
etcd_network_peer_sent_bytes_total{To="b10c280072f91b6"} 4.495762662e+09
etcd_network_peer_sent_bytes_total{To="ceb61f37051a244c"} 7.236551553e+09
# HELP etcd_network_peer_sent_failures_total The total number of send failures from peers.
# TYPE etcd_network_peer_sent_failures_total counter
etcd_network_peer_sent_failures_total{To="b10c280072f91b6"} 27
# HELP etcd_server_go_version Which Go version server is running with. 1 for 'server_go_version' label with current version.
# TYPE etcd_server_go_version gauge
etcd_server_go_version{server_go_version="go1.12.12"} 1
# HELP etcd_server_has_leader Whether or not a leader exists. 1 is existence, 0 is not.
# TYPE etcd_server_has_leader gauge
etcd_server_has_leader 1
# HELP etcd_server_health_failures The total number of failed health checks
# TYPE etcd_server_health_failures counter
etcd_server_health_failures 38
# HELP etcd_server_health_success The total number of successful health checks
# TYPE etcd_server_health_success counter
etcd_server_health_success 59829
# HELP etcd_server_heartbeat_send_failures_total The total number of leader heartbeat send failures (likely overloaded from slow disk).
# TYPE etcd_server_heartbeat_send_failures_total counter
etcd_server_heartbeat_send_failures_total 11669
# HELP etcd_server_id Server or member ID in hexadecimal format. 1 for 'server_id' label with current ID.
# TYPE etcd_server_id gauge
etcd_server_id{server_id="b3e330dca330e585"} 1
# HELP etcd_server_is_leader Whether or not this member is a leader. 1 if is, 0 otherwise.
# TYPE etcd_server_is_leader gauge
etcd_server_is_leader 0
# HELP etcd_server_is_learner Whether or not this member is a learner. 1 if is, 0 otherwise.
# TYPE etcd_server_is_learner gauge
etcd_server_is_learner 0
# HELP etcd_server_leader_changes_seen_total The number of leader changes seen.
# TYPE etcd_server_leader_changes_seen_total counter
etcd_server_leader_changes_seen_total 69
# HELP etcd_server_learner_promote_successes The total number of successful learner promotions while this member is leader.
# TYPE etcd_server_learner_promote_successes counter
etcd_server_learner_promote_successes 0
# HELP etcd_server_proposals_applied_total The total number of consensus proposals applied.
# TYPE etcd_server_proposals_applied_total gauge
etcd_server_proposals_applied_total 5.7080915e+07
# HELP etcd_server_proposals_committed_total The total number of consensus proposals committed.
# TYPE etcd_server_proposals_committed_total gauge
etcd_server_proposals_committed_total 5.7080915e+07
# HELP etcd_server_proposals_failed_total The total number of failed proposals seen.
# TYPE etcd_server_proposals_failed_total counter
etcd_server_proposals_failed_total 134
# HELP etcd_server_proposals_pending The current number of pending proposals to commit.
# TYPE etcd_server_proposals_pending gauge
etcd_server_proposals_pending 0
# HELP etcd_server_quota_backend_bytes Current backend storage quota size in bytes.
# TYPE etcd_server_quota_backend_bytes gauge
etcd_server_quota_backend_bytes 2.147483648e+09
# HELP etcd_server_read_indexes_failed_total The total number of failed read indexes seen.
# TYPE etcd_server_read_indexes_failed_total counter
etcd_server_read_indexes_failed_total 67
# HELP etcd_server_slow_apply_total The total number of slow apply requests (likely overloaded from slow disk).
# TYPE etcd_server_slow_apply_total counter
etcd_server_slow_apply_total 63091
# HELP etcd_server_slow_read_indexes_total The total number of pending read indexes not in sync with leader's or timed out read index requests.
# TYPE etcd_server_slow_read_indexes_total counter
etcd_server_slow_read_indexes_total 24
# HELP etcd_server_snapshot_apply_in_progress_total 1 if the server is applying the incoming snapshot. 0 if none.
# TYPE etcd_server_snapshot_apply_in_progress_total gauge
etcd_server_snapshot_apply_in_progress_total 0
# HELP etcd_server_version Which version is running. 1 for 'server_version' label with current version.
# TYPE etcd_server_version gauge
etcd_server_version{server_version="3.4.3"} 1
# HELP etcd_snap_db_fsync_duration_seconds The latency distributions of fsyncing .snap.db file
# TYPE etcd_snap_db_fsync_duration_seconds histogram
etcd_snap_db_fsync_duration_seconds_bucket{le="0.001"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.002"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.004"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.008"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.016"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.032"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.064"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.128"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.256"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="0.512"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="1.024"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="2.048"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="4.096"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="8.192"} 0
etcd_snap_db_fsync_duration_seconds_bucket{le="+Inf"} 0
etcd_snap_db_fsync_duration_seconds_sum 0
etcd_snap_db_fsync_duration_seconds_count 0
# HELP etcd_snap_db_save_total_duration_seconds The total latency distributions of v3 snapshot save
# TYPE etcd_snap_db_save_total_duration_seconds histogram
etcd_snap_db_save_total_duration_seconds_bucket{le="0.1"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="0.2"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="0.4"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="0.8"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="1.6"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="3.2"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="6.4"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="12.8"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="25.6"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="51.2"} 0
etcd_snap_db_save_total_duration_seconds_bucket{le="+Inf"} 0
etcd_snap_db_save_total_duration_seconds_sum 0
etcd_snap_db_save_total_duration_seconds_count 0
# HELP etcd_snap_fsync_duration_seconds The latency distributions of fsync called by snap.
# TYPE etcd_snap_fsync_duration_seconds histogram
etcd_snap_fsync_duration_seconds_bucket{le="0.001"} 0
etcd_snap_fsync_duration_seconds_bucket{le="0.002"} 0
etcd_snap_fsync_duration_seconds_bucket{le="0.004"} 0
etcd_snap_fsync_duration_seconds_bucket{le="0.008"} 190
etcd_snap_fsync_duration_seconds_bucket{le="0.016"} 345
etcd_snap_fsync_duration_seconds_bucket{le="0.032"} 368
etcd_snap_fsync_duration_seconds_bucket{le="0.064"} 373
etcd_snap_fsync_duration_seconds_bucket{le="0.128"} 376
etcd_snap_fsync_duration_seconds_bucket{le="0.256"} 380
etcd_snap_fsync_duration_seconds_bucket{le="0.512"} 380
etcd_snap_fsync_duration_seconds_bucket{le="1.024"} 380
etcd_snap_fsync_duration_seconds_bucket{le="2.048"} 380
etcd_snap_fsync_duration_seconds_bucket{le="4.096"} 380
etcd_snap_fsync_duration_seconds_bucket{le="8.192"} 380
etcd_snap_fsync_duration_seconds_bucket{le="+Inf"} 380
etcd_snap_fsync_duration_seconds_sum 4.558391028999997
etcd_snap_fsync_duration_seconds_count 380
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 1.9785e-05
go_gc_duration_seconds{quantile="0.25"} 3.9497e-05
go_gc_duration_seconds{quantile="0.5"} 5.4392e-05
go_gc_duration_seconds{quantile="0.75"} 8.8947e-05
go_gc_duration_seconds{quantile="1"} 0.066811763
go_gc_duration_seconds_sum 3.466559304
go_gc_duration_seconds_count 24829
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 1037
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.12.12"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 5.40655608e+08
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.442343722976e+12
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 3.046876e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 3.649213616e+09
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0.0003744391383233565
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 4.2825728e+07
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 5.40655608e+08
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 7.04888832e+08
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 5.6385536e+08
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 274088
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 4.70269952e+08
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 1.268744192e+09
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.602512560772257e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 3.649487704e+09
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 3472
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 1.86336e+06
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 3.588096e+06
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 8.12064112e+08
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.178396e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 6.324224e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 6.324224e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.325723896e+09
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 15
# HELP grpc_server_handled_total Total number of RPCs completed on the server, regardless of success or failure.
# TYPE grpc_server_handled_total counter
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Aborted",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="AlreadyExists",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 180
grpc_server_handled_total{grpc_code="Canceled",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="DataLoss",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="DeadlineExceeded",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="FailedPrecondition",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Internal",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="InvalidArgument",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="NotFound",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1073
grpc_server_handled_total{grpc_code="OK",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 8
grpc_server_handled_total{grpc_code="OK",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 67590
grpc_server_handled_total{grpc_code="OK",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 1
grpc_server_handled_total{grpc_code="OK",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 4.635505e+06
grpc_server_handled_total{grpc_code="OK",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 4
grpc_server_handled_total{grpc_code="OK",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1.46152e+06
grpc_server_handled_total{grpc_code="OK",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="OK",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="OutOfRange",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="PermissionDenied",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="ResourceExhausted",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unauthenticated",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 369
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 90
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 712
grpc_server_handled_total{grpc_code="Unavailable",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unimplemented",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 36
grpc_server_handled_total{grpc_code="Unknown",grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 4
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 0
grpc_server_handled_total{grpc_code="Unknown",grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
# HELP grpc_server_msg_received_total Total number of RPC stream messages received on the server.
# TYPE grpc_server_msg_received_total counter
grpc_server_msg_received_total{grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1073
grpc_server_msg_received_total{grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 8
grpc_server_msg_received_total{grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 67590
grpc_server_msg_received_total{grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_msg_received_total{grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 1
grpc_server_msg_received_total{grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 4.63591e+06
grpc_server_msg_received_total{grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_msg_received_total{grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 4
grpc_server_msg_received_total{grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1.461614e+06
grpc_server_msg_received_total{grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_received_total{grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 1007
grpc_server_msg_received_total{grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
# HELP grpc_server_msg_sent_total Total number of gRPC stream messages sent by the server.
# TYPE grpc_server_msg_sent_total counter
grpc_server_msg_sent_total{grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1073
grpc_server_msg_sent_total{grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 8
grpc_server_msg_sent_total{grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 67590
grpc_server_msg_sent_total{grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_msg_sent_total{grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 1
grpc_server_msg_sent_total{grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 4.635505e+06
grpc_server_msg_sent_total{grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_msg_sent_total{grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 4
grpc_server_msg_sent_total{grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1.46152e+06
grpc_server_msg_sent_total{grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_msg_sent_total{grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 4.203486e+06
grpc_server_msg_sent_total{grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
# HELP grpc_server_started_total Total number of RPCs started on the server.
# TYPE grpc_server_started_total counter
grpc_server_started_total{grpc_method="Alarm",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="AuthDisable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="AuthEnable",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Authenticate",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Check",grpc_service="grpc.health.v1.Health",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Compact",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1073
grpc_server_started_total{grpc_method="Defragment",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 8
grpc_server_started_total{grpc_method="DeleteRange",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Hash",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="HashKV",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="LeaseGrant",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 67590
grpc_server_started_total{grpc_method="LeaseKeepAlive",grpc_service="etcdserverpb.Lease",grpc_type="bidi_stream"} 0
grpc_server_started_total{grpc_method="LeaseLeases",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="LeaseRevoke",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="LeaseTimeToLive",grpc_service="etcdserverpb.Lease",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="MemberAdd",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="MemberList",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 1
grpc_server_started_total{grpc_method="MemberPromote",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="MemberRemove",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="MemberUpdate",grpc_service="etcdserverpb.Cluster",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="MoveLeader",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Put",grpc_service="etcdserverpb.KV",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Range",grpc_service="etcdserverpb.KV",grpc_type="unary"} 4.63591e+06
grpc_server_started_total{grpc_method="RoleAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="RoleDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="RoleGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="RoleGrantPermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="RoleList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="RoleRevokePermission",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Snapshot",grpc_service="etcdserverpb.Maintenance",grpc_type="server_stream"} 0
grpc_server_started_total{grpc_method="Status",grpc_service="etcdserverpb.Maintenance",grpc_type="unary"} 4
grpc_server_started_total{grpc_method="Txn",grpc_service="etcdserverpb.KV",grpc_type="unary"} 1.461614e+06
grpc_server_started_total{grpc_method="UserAdd",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="UserChangePassword",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="UserDelete",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="UserGet",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="UserGrantRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="UserList",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="UserRevokeRole",grpc_service="etcdserverpb.Auth",grpc_type="unary"} 0
grpc_server_started_total{grpc_method="Watch",grpc_service="etcdserverpb.Watch",grpc_type="bidi_stream"} 1005
grpc_server_started_total{grpc_method="Watch",grpc_service="grpc.health.v1.Health",grpc_type="server_stream"} 0
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 31740.47
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 160
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.068343296e+09
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.60191389639e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.65177344e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 25550
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
