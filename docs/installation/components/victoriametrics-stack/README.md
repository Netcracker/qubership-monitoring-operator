---
icon: lucide/database
---

# VictoriaMetrics Stack

VictoriaMetrics is the default time series backend. The VictoriaMetrics operator reconciles the custom resources
described below, and `PlatformMonitoring` configures each of them through its own section.

* [VictoriaMetrics Integration](victoriametrics.md) — send metrics to an external VictoriaMetrics installation.
* [VM Operator](vm-operator.md) — the operator that manages the custom resources on this page.
* [VMSingle](vmsingle.md) — single-node storage for metrics.
* [VMAgent](vmagent.md) — scrapes targets and writes to storage.
* [VMAlert](vmalert.md) — evaluates alerting and recording rules.
* [VMAlertManager](vmalertmanager.md) — routes and delivers alerts.
* [VMAuth](vmauth.md) — authentication proxy in front of the stack.
* [VMUser](vmuser.md) — credentials and access rules for VMAuth.
