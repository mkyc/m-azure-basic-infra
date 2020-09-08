# m-azure-basic-infra
Epiphany Module: Azure Basic Infrastructure

# Prepare service principal

Have a look [here](https://www.terraform.io/docs/providers/azurerm/guides/service_principal_client_secret.html).

```
az login 
az account list #get subscription from id field
az account set --subscription="SUBSCRIPTION_ID"
az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/SUBSCRIPTION_ID" --name="SOME_MEANINGFUL_NAME" #get appID, password, tenant, name and displayName
```

# Run module

```
cd examples/basic_flow
ARM_CLIENT_ID="appId field" ARM_CLIENT_SECRET="password field" ARM_SUBSCRIPTION_ID="id field" ARM_TENANT_ID="tenant field" make all
```

Or use config file with credentials:

```
cd examples/basic_flow
cat >azure.mk <<'EOF'
ARM_CLIENT_ID ?= "appId field"
ARM_CLIENT_SECRET ?= "password field"
ARM_SUBSCRIPTION_ID ?= "id field"
ARM_TENANT_ID ?= "tenant field"
EOF
make all
```
