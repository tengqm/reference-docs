---
title: Service Discovery and LoadBalancing
weight: 50
---

Discovery and load balancing resources are responsible for stitching your
workloads together into an accessible Loadbalanced Service.  By default,
[Workloads](../workloads/) are only accessible within the cluster, and they
must be exposed externally using a either a *LoadBalancer* or *NodePort*
[Service](../resources/service-v1-core/). For development, internally
accessible workloads can be accessed via proxy through the API server using
the `kubectl proxy` command.

Common resource types:

- [Services](../resources/service-v1-core/) for providing a single ip endpoint
  loadbalanced across multiple workload replicas.
- [Ingress](../resources/ingress-v1-network/) for providing a HTTP(S) endpoint
  routed to one or more *Services*

