---
title: Configuration and Storage
weight: 40
---

Config and Storage resources are responsible for injecting data into your
applications and persisting data externally to your container.

Common resource types:

- [ConfigMaps](../resources/configmap-v1-core/) for providing text key value
  pairs injected into the application through environment variables, command
  line arguments or files
- [Secrets](../resources/secret-v1-core/) for providing binary data injected
  into the application through files.
- [Volumes](../resources/volume-v1-core/) for providing a filesystem external
  to the container. Maybe shared across containers within the same Pod and
  have a lifetime persisting beyond a Container or Pod.

