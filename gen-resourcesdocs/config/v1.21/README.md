# Changes on v1.21

- `RollingUpdateDaemonSet`: new field `maxSurge`
- `EphemeralContainers`: added
- `NetworkPolicyPort`: new field `endport`
- `CSIStorageCapacity`: added
- `ConfigMap.immutable` GA
- `Secret.immutable` GA
- `PodSecurityPolicy` deprecated
- `ServiceSpec.topologyKeys` deprecated
- `CronJob` v2alpha1 removed
- `CronJob` GA
- `EndpointSlice` GA
- `JobSpec`: new fields `completionMode`, `suspend`
- `JobStatus`; new field `completedIndexes`
- `CronJobStatus`: new field `lastSuccessfulTime`
- `PodAffinityTerm`: new field `namespaceSelector`
- `ServiceSpec`: new fields `loadBalancerClass`, `internalTrafficPolicy`
- `IngressClassParametersReference` added
- `IngressClassSpec.parameters` of type `IngressClassParametersReference` (v1 and v1beta1)
- `PodDisruptionBudgetStatus`: new field `conditions`
- `EphemeralVolumeSource`: remove field `readonly`
- `Probe`: new field `terminationGracePeriodSeconds`
- `Endpoint`: new field `hints`
- new definition `EndpointHints`
- new definition `ForZone`
- `PodDisruptionBudget` GA
