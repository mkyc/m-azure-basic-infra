# m-azure-basic-infrastructure

Epiphany Module: Azure Basic Infrastructure

# Prepare service principal

Have a look [here](https://www.terraform.io/docs/providers/azurerm/guides/service_principal_client_secret.html).

```shell
az login 
az account list #get subscription from id field
az account set --subscription="SUBSCRIPTION_ID"
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/SUBSCRIPTION_ID" --name="SOME_MEANINGFUL_NAME" #get appID, password, tenant, name and displayName
```

# Build image

In main directory run:

```shell
make build
```

# Run module

```shell
cd examples/basic_flow
ARM_CLIENT_ID="appId field" ARM_CLIENT_SECRET="password field" ARM_SUBSCRIPTION_ID="id field" ARM_TENANT_ID="tenant field" make all
```

Or use config file with credentials:

```shell
cd examples/basic_flow
cat >azure.mk <<'EOF'
ARM_CLIENT_ID ?= "appId field"
ARM_CLIENT_SECRET ?= "password field"
ARM_SUBSCRIPTION_ID ?= "id field"
ARM_TENANT_ID ?= "tenant field"
EOF
make all
```

# Release module

```shell
make release
```

or if you want to set different version number:

```shell
make release VERSION=number_of_your_choice
```

# Run tests

```
make test
```

# Run tests in Kubernetes based build system

Kubernetes based build system means that build agents work inside Kubernetes cluster. During testing process application runs inside docker container. This means that we've got "docker inside docker (DiD)". This kind of environment requires a bit different configuration of mount shared storage to docker container than with standard one-layer configuration.

With DiD configuration shared volume needs to be created on host machine and this volume is shared with application container as Kubernetes volume.
Configuration steps:

1.  Create volume  (host path). In deployment.yaml add this config to create kubernetes volume:

```
volumes:
- name: tests-share
  hostPath:
    path: /tmp/tests-share
```

See manual for more details: https://kubernetes.io/docs/concepts/storage/volumes/#hostpath

2. Add mount point for kubernetes pod (agent). In deployment.yaml add this config to define volume's mount point:

```
volumeMounts:
- mountPath: /tests-share
  name: tests-share
```

3. Inside pod where tests will run set two variables to indicate host path and mount point:

```
export K8S_HOST_PATH=/tests-share
export K8S_VOL_PATH=/tmp/tests-share  ##modify paths according your needs, but they need to match paths from steps 1 and 2.
```

4. Go to location where you downloaded repository and run:

```
make test
```

5. Test results will be availabe inside ```/tests-share``` on pod on which tests are running and is mapped to ```/tmp/tests-share``` on kubernetes node.

# Windows users

This module is designed for Linux/Unix development/usage only. If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then ```examples/basic_flow_devcontainer``` folder.

# Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AzureRM provider | 2.27.0 | https://github.com/terraform-providers/terraform-provider-azurerm | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/LICENSE) |
| Make                      | 4.3     | https://www.gnu.org/software/make/                    | [ GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                        | 3.3.4   | https://github.com/mikefarah/yq/                      | [ MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
