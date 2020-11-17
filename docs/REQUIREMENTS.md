# Requirements

This document describes requirements to be able to run module and its tests. Note that latest software versions were used at the time of writing this document, so if some bug is found with using particular versions, it should be fixed and/or noted in this document.

## Service Principal

Creating a Service Principal using Azure CLI
Have a look [here](https://www.terraform.io/docs/providers/azurerm/guides/service_principal_client_secret.html).

```shell
az login
az account list #get subscription from id field
az account set --subscription="SUBSCRIPTION_ID"
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/SUBSCRIPTION_ID" --name="SOME_MEANINGFUL_NAME" #get appID, password, tenant, name and displayName
```

## Operating systems

The module is designed for Linux/Unix development/usage only.
If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then `examples/basic_flow_devcontainer` folder.

If there is a need to build the module container to check interactions between modules etc., Linux and MacOS should work fine and on Windows [WSL(2)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) can be used.
In WSL2 you can easily share the docker host between Windows and the WSL2 environment by selecting the WSL2 backend in Docker Desktop settings.
In WSL1 you can achieve the same with [this](https://nickjanetakis.com/blog/setting-up-docker-for-windows-and-wsl-to-work-flawlessly) tutorial.

## Build and release module

* Make
* Docker

## Run module

* Make - optional, required only if you want to run examples with Makefiles
* Docker

## Execute module tests

* Make
* Docker
* Golang > 1.15
