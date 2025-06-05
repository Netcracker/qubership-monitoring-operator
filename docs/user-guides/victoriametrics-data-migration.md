This guide describes how to migrate data between two VictoriaMetrics databases.

# Overview

VictoriaMetrics provides a utility called `vmctl` to migrate data between two VictoriaMetrics databases.
More information about internal working of `vmctl` can be found in the
[official documentation](https://docs.victoriametrics.com/vmctl/#migrating-data-from-victoriametrics).
`vmctl` provides `vm-native` option to perform data migration.

## Using Migration Job

While there are many ways to install and use `vmctl`, we shall make use of a Kubernetes job to perform data migration.
Below is an example job that can be used to migrate data from one Kubernetes cluster to another. In this particular
example, we are migrating data from a cluster containing `VMSingle` to a cluster containing `VMCluster`.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: vmctl-migrating-data
  namespace: monitoring
spec:
  template:
    spec:
      containers:
      - name: vmctl
        image: victoriametrics/vmctl:v1.104.0
        args:
          - vm-native
          - --vm-native-src-addr=https://vmsingle-monitoring.test.org
          - --vm-native-dst-addr=https://vminsert-k8s.monitoring:8480/insert/0/prometheus
          - --vm-native-filter-time-start=2022-01-01T00:00:00Z
          - --vm-native-src-insecure-skip-verify=true
          - --vm-native-dst-insecure-skip-verify=true
          - --vm-native-src-user=admin
          - --vm-native-src-password=admin
          - --vm-native-filter-match={__name__!~"vm_app.*"}
        resources:
          limits:
            cpu: '1'
            memory: 1Gi
          requests:
            cpu: 100m
            memory: 100Mi
      restartPolicy: Never
  backoffLimit: 2
```

Notice that the job uses a `vmctl` docker image and the `args` section contains the necessary options to
perform the data migration. The `--vm-native-src-addr` option specifies the source VictoriaMetrics
instance and the `--vm-native-dst-addr` option specifies the destination VictoriaMetrics instance.
Both Ingress URL and the internal service URL may be used for the source and destination addresses.
The `--vm-native-filter-time-start` option specifies the start time of the data to be migrated.

Migration speed can be adjusted via `--vm-concurrency` cmd-line flag, which controls the number of concurrent
workers busy with processing. Please note that each worker can load up to a single vCPU core on VictoriaMetrics.
So try to set it according to allocated CPU resources of your VictoriaMetrics destination installation.

`vmctl` supports `--vm-native-disable-http-keep-alive` to allow `vmctl` to use non-persistent HTTP connections to
avoid error use of closed network connection when running a heavy export requests.

**Note:**
In this case, we used  `--vm-native-src-insecure-skip-verify=true` and `--vm-native-dst-insecure-skip-verify=true`
to skip SSL verification for both source and destination VictoriaMetrics instances. This is not recommended.
Use `--vm-native-src-cert-file`, `--vm-native-src-key-file` and `--vm-native-src-ca-file` to specify TLS certificates
for source and `--vm-native-dst-cert-file` and `--vm-native-dst-key-file` to specify TLS certificates for destination.
An example is provided below.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: vmctl-migrating-data
  namespace: monitoring
spec:
  template:
    spec:
      containers:
      - name: vmctl
        image: victoriametrics/vmctl:v1.104.0
        args:
          - vm-native
          - --vm-native-src-addr=https://vmsingle-monitoring.test.org
          - --vm-native-dst-addr=https://vminsert-k8s.monitoring:8480/insert/0/prometheus
          - --vm-native-filter-time-start=2022-01-01T00:00:00Z
          - --vm-native-src-insecure-skip-verify=true
          - --vm-native-dst-insecure-skip-verify=true
          - --vm-native-src-user=admin
          - --vm-native-src-password=admin
          - --vm-native-filter-match={__name__!~"vm_app.*"}
          - --vm-native-src-cert-file=/etc/vm/secrets/src_tls.crt
          - --vm-native-src-key-file=/etc/vm/secrets/src_tls.key
          - --vm-native-dst-cert-file=/etc/vm/secrets/dst_tls.crt
          - --vm-native-dst-key-file=/etc/vm/secrets/dst_tls.key
        resources:
          limits:
            cpu: '1'
            memory: 1Gi
          requests:
            cpu: 100m
            memory: 100Mi
      restartPolicy: Never
  backoffLimit: 2
```

### Migration possibilities

Using the same job configuration and different parameters, it is possible to migrate data between VM instances:
single to single, cluster to cluster, single to cluster and vice versa. The source and destination instances
can be present in two different clusters or in the same cluster at different namespaces.

Below is the example of parameters that can be used in the job's `args` section to migrate data
between different types of VM instances:

```bash
# Migrating from cluster specific tenantID to single
--vm-native-src-addr=http://<src-vmselect>:8481/select/0/prometheus
--vm-native-dst-addr=http://<dst-vmsingle>:8428

# Migrating from single to cluster specific tenantID
--vm-native-src-addr=http://<src-vmsingle>:8428
--vm-native-dst-addr=http://<dst-vminsert>:8480/insert/0/prometheus

# Migrating single to single
--vm-native-src-addr=http://<src-vmsingle>:8428
--vm-native-dst-addr=http://<dst-vmsingle>:8428

# Migrating cluster to cluster for specific tenant ID
--vm-native-src-addr=http://<src-vmselect>:8481/select/0/prometheus
--vm-native-dst-addr=http://<dst-vminsert>:8480/insert/0/prometheus
```

### Job Results

The result of data migration can be found in the pod's (corresponding to the job) logs. Below is an example:

```bash
Requests to make: 13693 / 13694 [███████████████████████████████████████████████████████████] 99.99%
Requests to make: 13694 / 13694 [██████████████████████████████████████████████████████████] 100.00%
2024/10/14 14:21:07 Import finished!
2024/10/14 14:21:07 VictoriaMetrics importer stats:
  time spent while importing: 3h44m57.290974676s;
  total bytes: 6.5 GB;
  bytes/s: 483.1 kB;
  requests: 13694;
  requests retries: 1;
2024/10/14 14:21:07 Total time: 3h44m57.292391347s
```
