# Requirements

This document describes requirements to be able to run module. 

## Service Principal

Creating a Service Principal using Azure CLI
Have a look [here](https://www.terraform.io/docs/providers/azurerm/guides/service_principal_client_secret.html).

```shell
az login
az account list #get subscription from id field
az account set --subscription="SUBSCRIPTION_ID"
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/SUBSCRIPTION_ID" --name="SOME_MEANINGFUL_NAME" #get appID, password, tenant, name and displayName
```

## Run module

* Docker
* Make - optional, required only if you want to run examples with Makefiles
