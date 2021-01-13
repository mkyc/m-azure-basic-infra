# m-azure-basic-infrastructure

Epiphany Module: Azure Basic Infrastructure

AzBI module is responsible for providing basic cloud resources (eg. resource groups, virtual networks, subnets, virtual
machines etc.) which will be used by upcoming modules.

# Basic usage

## Requirements

Requirements are listed in a separate [document](docs/REQUIREMENTS.md).

## Run module

* Create a shared directory:

  ```shell
  mkdir /tmp/shared
  ```

  This 'shared' dir is a place where all configs and states will be stored while working with Epiphany modules.

* Generate ssh keys in: /tmp/shared/vms_rsa.pub

  ```shell
  ssh-keygen -t rsa -b 4096 -f /tmp/shared/vms_rsa -N ''
  ```

* Build Docker image if development version is used

  ```shell
  make build
  ```

* Initialize AzBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared epiphanyplatform/azbi:dev init --name=epiphany-modules-test
  ```

  :star: Variable values can be passed as docker environment variables as well. In presented example we could
  use `docker run` command `-e NAME=2` parameter instead of `--name=epiphany-modules-test` command parameter.

  :warning: Use image's tag according to tag generated in build step.

  This command will create configuration file of AzBI module in /tmp/shared/azbi/azbi-config.yml. You can investigate
  what is stored in that file. Available parameters are described in the [inputs](docs/INPUTS.adoc) document.

  :warning: Pay attention to the docker image tag you are using. Command `make build` command uses a specific version
  tag (default `epiphanyplatrofm/azbi:dev`).

* Plan and apply AzBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azbi:dev plan
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azbi:dev apply
  ```
  :star: Variable values can be passed as docker environment variables. I's often more convenient to pass sensitive
  values as presented.

  Running those commands should create resource group, vnet, subnet and 2 virtual machines. You should verify in Azure
  Portal.

* Destroy module resources:

  ```shell
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azbi:dev plan --destroy
  docker run --rm -v /tmp/shared:/shared -e SUBSCRIPTION_ID=subscriptionId -e CLIENT_ID=appId -e CLIENT_SECRET=password -e TENANT_ID=tenantId epiphanyplatform/azbi:dev destroy
  ```
  :star: Variable values can be passed as docker environment variables. I's often more convenient to pass sensitive values as presented.
  
  :warning: Running those commands will remove all resource group and all its content so be careful. You should verify in Azure Portal.

# AzBI output data

The output from this module is:

* rg_name
* vnet_name
* vm_groups

# Examples

For examples running description please have a look into [this document](docs/EXAMPLES.md).

# Development

For development related topics please look into [this document](docs/DEVELOPMENT.md).

# Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AzureRM provider | 2.27.0 | https://github.com/terraform-providers/terraform-provider-azurerm | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/LICENSE) |
