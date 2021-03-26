---
title: Cluster Resources
weight: 70
---

Cluster resources are responsible for defining configuration of the cluster
itself, and are generally only used by cluster operators.

Example resource types:

- [`ResourceQuota`](../resources/resourcequota-v1-core/) defines the quota of
  resource usage by Pods and containers.
- [`ServiceAccount`](../resources/serviceaccount-v1-core/) defines the service
  account for accessing cluster resources.
