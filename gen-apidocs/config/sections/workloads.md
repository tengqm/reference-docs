---
title: Workload Resources
weight: 30
---

Workloads resources are responsible for managing and running your containers
on the cluster. [Containers](#container-v1-core) are created by controllers
through [Pods](#pod-v1-core). Pods run containers and provide environmental
dependencies such as shared or persistent storage [Volumes](#volume-v1-core)
and [ConfigMap](#configmap-v1-core) or [Secret](#secret-v1-core) data
injected into the container.

The most common Controllers are:

- [Deployments](../resources/deployment-v1-apps/) for stateless persistent
  applications such as HTTP servers.
- [StatefulSets](../resources/statefulset-v1-apps/) for stateful persistent
  applications such as database servers.
- [Jobs](../resources/job-v1-batch/) for run-to-completion applications such
  applications such as batch jobs.
- [DaemonSets](../resources/daemonset-v1-apps/) for background tasks that run
  on cluster nodes for cluster operations.

