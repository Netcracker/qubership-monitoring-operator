# HA Services

This dashboard allows checking whether services in the cluster are running in high availability mode or not

## Tags

* `ha`

## Panels

### Help

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| How to | Describes how to use this dashboard |  |  |
<!-- markdownlint-enable line-length -->

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Services with Available Replicas | Shows total number of services (Deployments and StatefulSets) with at least 1 available replica, number of services with only 1 available replica (non-HA mode) and with 2 or more replicas (HA mode) |  |  |
| Services with Unavailable Replicas | Total number of services (Deployments and StatefulSets) with unavailable replicas | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Services with Multiple Pods per Node | Total number of services (Deployments and StatefulSets) with 2 or more pods placed on the same node | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Total Number of Replicas | Shows number of pods controlled by Deployments and StatefulSets |  |  |
<!-- markdownlint-enable line-length -->

### Deployments

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Desired and Available Replicas | Shows number of available and desired replicas for every Deployment<br/><br/>Services in HA mode must have more than 1 available replicas |  |  |
| Multiple Replicas on the Same Node | Shows Deployments if 2 or more their replicas placed on the same node<br/><br/>Services are not allowed to have 2 or more replicas on the same node in HA mode |  |  |
| Not Healthy Pods | Show information about the reason the container is currently in waiting or terminated state for pods controlled by Deployments | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Terminated Status | Show information about the last reason the container was in terminated state for pods controlled by Deployments | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### StatefulSets

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Desired and Available Replicas | Shows number of available and desired replicas for every StatefulSet<br/><br/>Services in HA mode must have more than 1 available replicas |  |  |
| Multiple Replicas on the Same Node | Shows StatefulSets if 2 or more their replicas placed on the same node<br/><br/>Services are not allowed to have 2 or more replicas on the same node in HA mode |  |  |
| Not Healthy Pods | Show information about the reason the container is currently in waiting or terminated state for pods controlled by StatefulSets | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Terminated Status | Show information about the last reason the container was in terminated state for pods controlled by StatefulSets | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
