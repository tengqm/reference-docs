---
title: Metadata Objects
weight: 60
---

Metadata resources are responsible for configuring behavior of your other
resources within the cluster.

Common resource types:

- [HorizontalPodAutoscaler](../resources/horizontalpodautoscaler-v1-autoscaling/)
  (HPA) for automatically scaling the replicacount of your workloads in
  response to load variations.
- [PodDisruptionBudget](../resources/poddisruptionbudget-v1beta1-policy) for
  configuring how many replicas in a given workload maybe made concurrently
  unavailable when performing maintenance.
- [CustomResourceDefinitions](../resources/customresourcedefinitions-v1-apiextensions)
  for extending the Kubernetes APIs with your own types.
- [Event](../resources/event-v1-core/) for notification of resource lifecycle
  events in the cluster.

