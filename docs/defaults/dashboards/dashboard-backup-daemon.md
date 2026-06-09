# Backup Daemon

Backup Information Dashboard.

## Tags

* `backups`

## Panels

### Backup Daemon Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Status | Shows Backup Adapter status | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Status Transition | Shows Backup Adapter status transition |  |  |
| CPU Usage | Shows CPU usage per container in pod against resource requests and limits |  |  |
| Memory Usage | Shows memory usage per container in pod against resource requests and limits |  |  |
| Last Backup Status | Shows Last Backup status | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Last Backup Status Transition | Shows Last Backup status transition |  |  |
| Last Backup Size | Size of last successful backup  |  |  |
| Dumps count | Current number of dumps |  |  |
| Storage Space Usage | Shows used by backups space against total space |  |  |
| Storage Inodes Usage | Shows inodes usage |  |  |
<!-- markdownlint-enable line-length -->
