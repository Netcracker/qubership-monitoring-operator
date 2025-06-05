# Prerequisites

This section describes the prerequisites that should be checked before deploying the Qubership Monitoring Operator.

## Before You Begin

Verify the following requirements before installation:

### System Requirements

* **Kubernetes**: Version 1.19+ or OpenShift 3.11+
* **CLI Tools**: `kubectl` 1.19+ or `oc` 3.11+
* **Package Manager**: Helm 3.0+
* **Namespace**: Target namespace must be created before installation
* **Permissions**: Sufficient rights to deploy in the target project (see [Permissions](#permissions))

### Important Notes

* If you disable CRD upgrade automation, you need to apply ALL CRDs before installation
* If updating an existing deployment, read the [Maintenance Guide](../maintenance.md) first
* Pod Security Policy restrictions may affect `node-exporter` pods - add privileged PSP ClusterRole if needed
* If other Prometheus/VictoriaMetrics instances exist, check namespace scope restrictions
* For OpenShift deployments, verify the namespace has an empty `node-selector`

Example OpenShift namespace check:
```bash
$ oc describe namespace <namespace>
Name:         <namespace>
Labels:       <none>
Annotations:  openshift.io/description=
              openshift.io/display-name=
              openshift.io/node-selector=
              openshift.io/requester=admin
Status:        Active
```

## Hardware Requirements

### Minimal Requirements

The **minimal** recommended hardware requirements for an average system:

| Resource | Requirement |
| -------- | ----------- |
| CPU      | 3000m       |
| RAM      | 3000Mi      |
| Storage  | 10GB        |

### Component Resource Breakdown

Default resource allocation by component:

| Component             | Requests                | Limits                  |
| --------------------- | ----------------------- | ----------------------- |
| `alertmanager`        | CPU: 100m, RAM: 100Mi   | CPU: 200m, RAM: 200Mi   |
| `grafana`             | CPU: 300m, RAM: 400Mi   | CPU: 500m, RAM: 800Mi   |
| `grafana-operator`    | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 100Mi   |
| `kube-state-metrics`  | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 256Mi   |
| `monitoring-operator` | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 150Mi   |
| `node-exporter`       | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 100Mi   |
| `victoriametrics`     | CPU: 1000m, RAM: 3Gi    | CPU: 1500m, RAM: 5Gi    |
| `vm-operator`         | CPU: 100m, RAM: 100Mi   | CPU: 200m, RAM: 200Mi   |
| **Total**             | CPU: 1700m, RAM: 3872Mi | CPU: 2800m, RAM: 6926Mi |

!!! note "Resource Sizing"
    We strongly recommend reviewing and adjusting these requests and limits based on your environment's specific needs.

For detailed resource calculations, see the [Resource Usage](../resource-usage.md) guide.

## Permissions

The `monitoring-operator` service account requires ClusterRole permissions for various Kubernetes APIs.

### Required ClusterRoles

Components that require ClusterRole creation:

* `prometheus` / `prometheus-operator`
* `grafana-operator`
* `kube-state-metrics`
* `node-exporter`
* `cert-exporter`
* `cloudwatch-exporter`
* `network-latency-exporter`
* `prometheus-adapter-operator`
* `version-exporter`
* `vmoperator`

### Restricted Deployment

The monitoring-operator supports deployment with restricted privileges:
- Access limited to resources within the namespace
- Restricted access to cluster-scoped resources
- Manual privilege grants required for cluster-scoped resources and other namespaces

## Cloud Platform Specifics

### AWS

#### EBS Persistent Volumes (Recommended)

For AWS deployments, use **AWS EBS** for persistent volumes:

- **Prometheus**: SSD volumes recommended
- **VictoriaMetrics**: Minimum HDD with "Throughput Optimized HDD" type

See [AWS EBS Volume Types](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html) for details.

#### EFS Persistent Volumes

!!! warning "Prometheus EFS Compatibility"
    **Prometheus does NOT support** non-POSIX compliant storages like AWS EFS. Avoid EFS for Prometheus deployments.

VictoriaMetrics supports EFS but requires proper throughput planning:
- Avoid "Bursting" throughput mode
- Use "Provisioned Throughput" (minimum 10 MB/s) or "Elastic Throughput"

## Supported Versions

### Kubernetes Compatibility

Current recommended Kubernetes version: **1.26.x**

| Kubernetes Version     | Status           |
| ---------------------- | ---------------- |
| `1.24.x`               | Tested           |
| `1.25.x`               | Tested           |
| `1.26.x` (recommended) | Tested           |
| `1.27.x`               | Forward compatible |
| `1.28.x`               | Forward compatible |

### Cloud Platform Support

| Platform                                 | Support |
| ---------------------------------------- | ------- |
| AWS Elastic Kubernetes Service (AWS EKS) | ✓       |
| Azure Kubernetes Service (AKS)          | ✓       |
| Google Kubernetes Engine (GKE)          | ✓       |
| On-premise Kubernetes >= 1.25           | ✓       |
| On-premise OpenShift >= 4.10            | ✓       |

### Legacy Version Support

For older Kubernetes/OpenShift versions:

| Kubernetes  | Last Supported Monitoring Version |
| ----------- | --------------------------------- |
| `< v1.25.0` | `0.46.0`                         |
| `< v1.18.0` | `0.26.0`                         |

| OpenShift | Last Supported Monitoring Version |
| --------- | --------------------------------- |
| `v3.11.0` | `0.46.0`                         |

## Next Steps

After verifying prerequisites:

1. Review [Basic Components](basic-components.md) to understand what will be installed
2. Configure [Storage](storage.md) for persistent data
3. Proceed with [Deployment](deploy.md) 