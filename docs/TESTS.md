# Tests

## Prepare service principal

Prepare service principal variables file before running this tests:

```shell
CLIENT_ID=xxx CLIENT_SECRET=yyy SUBSCRIPTION_ID=zzz TENANT_ID=vvv make prepare-service-principal
```

## Run module tests

In the main directory, run:

```shell
make test
```

## Run tests in Kubernetes based build system

Kubernetes based build system means that build agents work inside Kubernetes cluster. During testing process application runs inside docker container. This means that we've got "docker inside docker (DiD)". This kind of environment requires a bit different configuration of mount shared storage to docker container than with standard one-layer configuration.

With DiD configuration shared volume needs to be created on host machine and this volume is shared with application container as Kubernetes volume.

Configuration steps:

1. Create volume (host path). In deployment.yaml add this config to create kubernetes volume:

```yaml
volumes:
- name: tests-share
  hostPath:
    path: /tmp/tests-share
```

See manual for more details: https://kubernetes.io/docs/concepts/storage/volumes/#hostpath

2. Add mount point for kubernetes pod (agent). In deployment.yaml add this config to define volume's mount point:

```yaml
volumeMounts:
- mountPath: /tests-share
  name: tests-share
```

3. Inside pod where tests will run set two variables to indicate host path and mount point:

```shell
export K8S_HOST_PATH=/tests-share
export K8S_VOL_PATH=/tmp/tests-share  ##modify paths according your needs, but they need to match paths from steps 1 and 2.
```

4. Go to location where you downloaded repository and run:

```shell
make test
```

5. Test results will be availabe inside ```/tests-share``` on pod on which tests are running and is mapped to ```/tmp/tests-share``` on kubernetes node.
