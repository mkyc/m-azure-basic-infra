# m-azure-basic-infrastructure

Epiphany Module: Azure Basic Infrastructure

AzBI module is responsible for providing basic cloud resources (eg. resource groups, virtual networks, subnets, virtual machines ect.) which will be used by upcoming modules.

# Basic usage

## Requirements

Requirements are listed in a separate [document](docs/REQUIREMENTS.md).

## Build image

In main directory, run:

  ```shell
  make build
  ```

  Note: This command uses the default VERSION variable.

or directly using Docker:

  ```shell
  cd m-azure-basic-infrastructure/
  docker build --tag epiphanyplatform/azbi:latest .
  ```

Note: Re-run the above commands will overwrite your docker image. To bypass that, specify a different tag name or VERSION variable.

## Run module

* Create a shared directory:

  ```shell
  mkdir /tmp/shared
  ```

  This 'shared' dir is a place where all configs and states will be stored while working with modules.

* Generate ssh keys in: /tmp/shared/vms_rsa.pub

  ```shell
  ssh-keygen -t rsa -b 4096 -f /tmp/shared/vms_rsa -N ''
  ```

* Initialize AzBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azbi:latest init M_VMS_COUNT=2 M_PUBLIC_IPS=true M_NAME=epiphany-modules-test
  ```

  This command will create configuration file of AzBI module in /tmp/shared/azbi/azbi-config.yml. You can investigate what is stored in that file. Available parameters are listed in the [inputs](docs/INPUTS.adoc) document.

  Note: Pay attention to the docker image tag you are using. `make build` command uses a specific version tag eg. epiphanyplatrofm/azbi:0.0.1.

* Plan and apply AzBI module:

  ```shell
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azbi:latest plan M_ARM_CLIENT_ID=appId M_ARM_CLIENT_SECRET=password M_ARM_SUBSCRIPTION_ID=subscriptionId M_ARM_TENANT_ID=tenantId
  docker run --rm -v /tmp/shared:/shared -t epiphanyplatform/azbi:latest apply M_ARM_CLIENT_ID=appId M_ARM_CLIENT_SECRET=password M_ARM_SUBSCRIPTION_ID=subscriptionId M_ARM_TENANT_ID=tenantId
  ```

  Running those commands should create resource group, vnet, subnet and 2 virtual machines. You should verify in Azure Portal.

## Run module with provided example

### Prepare config file

Prepare your own variables in azure.mk file to use in the building process.
Sample file (examples/basic_flow/azure.mk.sample):

  ```shell
  ARM_CLIENT_ID ?= "appId field"
  ARM_CLIENT_SECRET ?= "password field"
  ARM_SUBSCRIPTION_ID ?= "id field"
  ARM_TENANT_ID ?= "tenant field"
  ```

* Create an environment

  ```shell
  cd examples/basic_flow
  make all
  ```

* Delete an environment

  ```shell
  make destroy-plan
  make destroy
  ```

## Release module

  ```shell
  make release
  ```

or if you want to set different version number:

  ```shell
  make release VERSION=number_of_your_choice
  ```

## Run tests

Tests are described in a separate [document](docs/TESTS.md).

# AzBI output data

The output from this module is:

* private_ips
* public_ips
* vm_names
* rg_name
* vnet_name

# Module dependencies

| Component                 | Version | Repo/Website                                          | License                                                           |
| ------------------------- | ------- | ----------------------------------------------------- | ----------------------------------------------------------------- |
| Terraform                 | 0.13.2  | https://www.terraform.io/                             | [Mozilla Public License 2.0](https://github.com/hashicorp/terraform/blob/master/LICENSE) |
| Terraform AzureRM provider | 2.27.0 | https://github.com/terraform-providers/terraform-provider-azurerm | [Mozilla Public License 2.0](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/LICENSE) |
| Make                      | 4.3     | https://www.gnu.org/software/make/                    | [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html) |
| yq                        | 3.3.4   | https://github.com/mikefarah/yq/                      | [MIT License](https://github.com/mikefarah/yq/blob/master/LICENSE) |
