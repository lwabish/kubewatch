# customization
Based on kubewatch, the feature-permission branch adds a handler called permission.

This handler listens all new pvs and filter them with specific storage class,after which the handler uses a kubernetes job
mounting the pvc and change the volume's permission and owner based on the handler's config.

The permission handler resolves this problem: root user is not allowed in docker image,but volume mounted is owned by root user.

