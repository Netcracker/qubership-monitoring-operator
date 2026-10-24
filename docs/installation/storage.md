# Storage Configuration

To maintain data across deployments and version upgrades, the data must be persisted to some volume
other than `emptyDir`, allowing it to be reused by Pods after an upgrade.

!!! note "Version Requirement"
    Ability to specify storage was added since version `0.4.0`

## Overview

Kubernetes supports several kinds of storage volumes. The Monitoring Operator uses Prometheus/VictoriaMetrics Operator
to control the Prometheus/VictoriaMetrics deployment. The Prometheus/VictoriaMetrics Operator works
with `PersistentVolumeClaims`, which supports the underlying `PersistentVolume` to be provisioned when requested.

This document assumes a basic understanding of `PersistentVolumes`, `PersistentVolumeClaims`, and their provisioning.

## Storage Types

### Dynamic Provisioning (Recommended)

Automatic provisioning of storage requires an already existing `StorageClass`.
For best results, use volumes that have high I/O throughput.

#### Prometheus Storage

```yaml
prometheus:
  storage:
    volumeClaimTemplate:
      spec:
        # Specify storage class to create volume
        storageClassName: nfs-dynamic-provisioning
        resources:
          requests:
            storage: 10Gi
```

#### VictoriaMetrics Storage

```yaml
victoriametrics:
  vmSingle:
    storage:
      # Specify storage class to create volume
      storageClassName: nfs-dynamic-provisioning
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      selector:
        matchLabels:
          app.kubernetes.io/name: vmsingle
      volumeName: data-victoriametrics-pv
```

### Manual Storage Provisioning

The monitoring deploy parameters allow you to support arbitrary storage through a PersistentVolumeClaim.
The easiest way to use a volume that cannot be automatically provisioned is to use a label selector
alongside a manually created PersistentVolume.

#### Prometheus Manual Storage

```yaml
prometheus:
  # Vanilla Prometheus image use user and group nobody = 65534
  # So for use PV it is better to use user nobody = 65534
  securityContext:
    fsGroup: 65534
    runAsUser: 65534
  # Because hostPath PV created on specific node, we must bind Prometheus on this node
  nodeSelector:
    kubernetes.io/hostname: worker1
  storage:
    volumeClaimTemplate:
      spec:
        resources:
          requests:
            storage: 10Gi
        selector:
          # Match PV by label on PV
          matchLabels:
            app: prometheus
```

#### VictoriaMetrics Manual Storage

```yaml
victoriametrics:
  vmSingle:
    storage:
      storageClassName: manual
      resources:
        requests:
          storage: 10Gi
      selector:
        # Match PV by label on PV
        matchLabels:
          app: victoriametrics
```

## Known Issues

### Prometheus Operator Version Issues

If you deploy monitoring with Prometheus Operator version `v0.34.0` or less, you may encounter:

- Creation/Update of Prometheus object throws error on volume claim template creationTimeStamp being null ([Issue #2824](https://github.com/prometheus-operator/prometheus-operator/issues/2824))
- Setting podMetadata labels without creation timestamp fails ([Issue #2399](https://github.com/prometheus-operator/prometheus-operator/issues/2399))

**Solution**: This issue was fixed since version `v0.35.0`. For older versions, manually add `creationTimeStamp`:

```yaml
prometheus:
  storage:
    volumeClaimTemplate:
      metadata:
        creationTimestamp: '2020-06-23T18:55:02Z'
      spec:
        # ... rest of configuration
```

### HostPath Permission Issues

When using `hostPath` with non-default users (specified with `securityContext`), you may encounter permission denied errors:

- Permission denied writing to mount using volumeClaimTemplate ([Issue #966](https://github.com/prometheus-operator/prometheus-operator/issues/966))
- mkdir /prometheus/wal: permission denied ([Issue #12176](https://github.com/helm/charts/issues/12176))

**Root Cause**: The Prometheus Docker image:

1. Creates a directory for Prometheus data under `root`
2. Changes ownership to user `nobody` (uid:gid = 65534:65534)
3. Runs under user `nobody`

**Solution**: If you need to run with a custom user, create `prometheus-db` in PV folder and set proper ownership:

```bash
# Create directory
mkdir -p /mnt/data/prometheus/prometheus-db

# Set ownership to custom user (example: uid:gid = 2001:2001)
chown -R 2001:2001 /mnt/data/prometheus
```

## Storage Classes

### Disabling Default StorageClasses

To manually provision volumes (Kubernetes 1.6.0+), you may need to disable the default StorageClass
that is automatically created for certain Cloud Providers.

The default StorageClass behavior overrides manual storage provisioning, preventing PersistentVolumeClaims
from automatically binding to manually created PersistentVolumes.

**Example**: Disable default StorageClass on minikube:

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  namespace: kube-system
  name: standard
  annotations:
    # disable this default storage class by setting this annotation to false.
    storageclass.beta.kubernetes.io/is-default-class: "false"
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
provisioner: k8s.io/minikube-hostpath
```

## Best Practices

### Storage Performance

1. **Use SSDs**: For Prometheus, use SSD storage for better IOPS
2. **Size Appropriately**: Start with 10-50GB, monitor usage patterns
3. **Backup Strategy**: Implement regular backups of storage volumes

### Cloud-Specific Recommendations

#### AWS

- **Prometheus**: Use gp3/io1 EBS volumes
- **VictoriaMetrics**: gp2/gp3 for cost-effective storage

#### Azure

- **Prometheus**: Premium SSD for production workloads
- **VictoriaMetrics**: Standard SSD for balanced performance

#### GCP

- **Prometheus**: SSD persistent disks
- **VictoriaMetrics**: Standard persistent disks acceptable

### Grafana Storage

```yaml
grafana:
  dataStorage:
    class: fast-ssd
    size: 10Gi
    accessModes:
      - ReadWriteOnce
```

### Common PVC Annotations

Use `pvc.metadata.annotations` to add the same annotations, such as `argocd.argoproj.io/sync-options: Prune=false`,
to the PVCs of the components listed below. For example:

```yaml
pvc:
  metadata:
    annotations:
      argocd.argoproj.io/sync-options: Prune=false
```

The chart merges these annotations into the PVC metadata of the following components:

| Component      | Component-specific parameter                                                      |
| -------------- | --------------------------------------------------------------------------------- |
| Prometheus     | `prometheus.storage.volumeClaimTemplate.metadata.annotations`                     |
| VMSingle       | `victoriametrics.vmSingle.storageMetadata.annotations`                            |
| VMAlertManager | `victoriametrics.vmAlertManager.storage.volumeClaimTemplate.metadata.annotations` |
| Grafana        | `grafana.dataStorage.annotations`                                                 |

Rules:

- A component-specific annotation overrides a common annotation with the same key.
- The chart adds annotations only to components with persistent storage configured: `storage.volumeClaimTemplate` for
  Prometheus and VMAlertManager, `storage` for VMSingle, and `dataStorage` for Grafana. The parameter never enables
  storage on its own.
- Pushgateway and VMCluster do not use this parameter.
- Annotation values must be strings. Quote values such as `"true"` in YAML, and use `--set-string` on the command line:

  ```bash
  helm upgrade ... --set-string 'pvc.metadata.annotations.argocd\.argoproj\.io/sync-options=Prune=false'
  ```

Limitations:

- Operators create these PVCs, not Helm. VMSingle and Grafana PVCs have an owner chain that leads to the
  PlatformMonitoring resource. Prometheus and VMAlertManager PVCs come from StatefulSet `volumeClaimTemplates` and
  carry labels with the custom resource name, not the Argo CD Application name. Argo CD therefore does not prune
  these PVCs with either `annotation` or `label` resource tracking. The annotations keep PVC metadata consistent and
  serve external tools that read them.
- Prometheus creates PVCs from StatefulSet `volumeClaimTemplates`. Kubernetes does not update an existing PVC when the
  template changes, so an upgrade does not annotate Prometheus PVCs that already exist. Changing the template also
  makes the operator recreate the Prometheus StatefulSet, which restarts Prometheus. Prometheus PVCs are named
  `prometheus-k8s-db-prometheus-k8s-<N>`, one per replica. List them with `kubectl get pvc -n <namespace>` and
  annotate each one manually:

  ```bash
  kubectl annotate pvc --overwrite -n <namespace> prometheus-k8s-db-prometheus-k8s-0 \
    argocd.argoproj.io/sync-options=Prune=false
  ```

- VMAlertManager PVCs keep annotation keys that you remove from the parameter. Remove them with
  `kubectl annotate pvc -n <namespace> <pvc-name> <key>-`.
- The annotations do not stop the operator from deleting the Prometheus PVC when you set `prometheus.install: false`
  or remove Prometheus storage. They also do not stop Kubernetes garbage collection from deleting PVCs whose owner
  chain leads to the PlatformMonitoring resource (VMSingle, Grafana) when that resource is deleted.

## Storage Examples

### Complete Production Setup

```yaml
prometheus:
  storage:
    volumeClaimTemplate:
      spec:
        storageClassName: fast-ssd
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 50Gi

victoriametrics:
  vmSingle:
    storage:
      storageClassName: fast-ssd
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 100Gi

grafana:
  dataStorage:
    class: fast-ssd
    size: 10Gi
    accessModes:
      - ReadWriteOnce
```

### Development Setup (hostPath)

```yaml
prometheus:
  securityContext:
    fsGroup: 65534
    runAsUser: 65534
  nodeSelector:
    kubernetes.io/hostname: dev-node
  storage:
    volumeClaimTemplate:
      spec:
        storageClassName: manual
        resources:
          requests:
            storage: 5Gi
        selector:
          matchLabels:
            app: prometheus-dev

victoriametrics:
  vmSingle:
    storage:
      storageClassName: manual
      resources:
        requests:
          storage: 10Gi
      selector:
        matchLabels:
          app: victoriametrics-dev
```

## Related Configuration

- **[Prerequisites](prerequisites.md)** - Hardware and storage requirements
- **[Basic Components](basic-components.md)** - Component resource requirements
- **[VictoriaMetrics Configuration](components/victoriametrics-stack/victoriametrics.md)** - VictoriaMetrics-specific
  storage options
