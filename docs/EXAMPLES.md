# Examples

## Prepare credentials

Prepare your own variables in `azure.mk` file to use in the process.
Sample file (`examples/basic_flow/azure.mk.sample`):

```shell
ARM_CLIENT_ID ?= client-id-value
ARM_CLIENT_SECRET ?= client-secret-field
ARM_SUBSCRIPTION_ID ?= subscription-id-field
ARM_TENANT_ID ?= tenant-id-field
```

# Create cluster

```shell
cd examples/basic_flow
make all
```

# Delete cluster

```shell
make destroy-plan
make destroy
```
