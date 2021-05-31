# Development

This document describes requirements to be able to develop abd build module. 

## Operating systems

We develop this module in Linux/MacOs environments mostly. 

If you need to develop from Windows you can use the included [devcontainer setup for VScode](https://code.visualstudio.com/docs/remote/containers-tutorial) and run the examples the same way but then from then `examples/basic_flow_devcontainer` folder.

To build module container Linux and MacOS should work fine. 

On Windows [WSL(2)](https://docs.microsoft.com/en-us/windows/wsl/install-win10) can be used. In WSL2 you can easily share the docker host between Windows and the WSL2 environment by selecting the WSL2 backend in Docker Desktop settings. In WSL1 you can achieve the same with [this](https://nickjanetakis.com/blog/setting-up-docker-for-windows-and-wsl-to-work-flawlessly) tutorial.


# Build image

In the main directory, run:

```shell
make build
```

:warning: This command uses the default VERSION variable (default: `dev`).

or directly using Docker:

```shell
go mod vendor
docker build --tag epiphanyplatform/azbi:dev .
```

:warning: Re-run the above commands will overwrite your existing docker image (if exists). To bypass that, specify a different `--tag` parameter. 

# Release module

```shell
make release
```

or if you want to set different version number:

```shell
make release VERSION=something_else
```

# Run tests

Tests are described in a separate [document](TESTS.md).

# Develop locally with e-structures repository

Assuming you work with [e-structures](https://github.com/epiphany-platform/e-structures) repository on branch named `some-branch-name`, you should run: 


```shell
go get github.com/epiphany-platform/e-structures@some-branch-name
go mod vendor
```

this would update `go.mod` file and download `e-structures` version to `vendor` directory. 

# Develop terraform scripts

To develop terraform scripts independently of go module code.

1) have image built already
1) go to `resources` directory in terminal
1) run `mkdir shared` in that `resources` directory
1) run `ssh-keygen -t rsa -b 4096 -f ./shared/vms_rsa -N ''` to generate required key pair
1) run module init command ie.: `docker run --rm -v $(pwd)/shared:/shared -t epiphanyplatform/azbi:dev init`
1) run module plan command with debug switch ie.: `docker run --rm -v $(pwd)/shared:/shared -e SUBSCRIPTION_ID=xxx -e CLIENT_ID=yyy -e CLIENT_SECRET=zzz -e TENANT_ID=vvv epiphanyplatform/azbi:dev plan --debug`
1) that would provide you detailed output and one of first lines provides used .tfvars.json file. It will look similar to following: 
   ```
   ...
   2021-05-31T11:38:41Z DBG ../go/src/github.com/epiphany-platform/m-azure-basic-infrastructure/cmd/common.go:61 > templateTfVars
   2021-05-31T11:38:41Z INF ../go/src/github.com/epiphany-platform/m-azure-basic-infrastructure/cmd/common.go:68 > {"name":"epiphany","location":"northeurope","address_space":["10.0.0.0/16"],"subnets":[{"name":"main","address_prefixes":["10.0.1.0/24"]}],"vm_groups":[{"name":"vm-group0","vm_count":1,"vm_size":"Standard_DS2_v2","use_public_ip":true,"subnet_names":["main"],"vm_image":{"publisher":"Canonical","offer":"UbuntuServer","sku":"18.04-LTS","version":"18.04.202006101"},"data_disks":[{"disk_size_gb":10,"storage_type":"Premium_LRS"}]}],"rsa_pub_path":"/shared/vms_rsa.pub"}
   ...
   ```
1) you should copy provided JSON to file `./terraform/terraform.tfvars.json` relatively to `resources` directory
1) now you can run terraform directly from `resources` directory

To be able to run version of terraform that scripts require you can use following docker snippet: 

```shell
docker run --rm -v $(pwd)/terraform:/workspace -v $(pwd)/shared:/shared -w /workspace hashicorp/terraform:0.13.7 init
docker run --rm -v $(pwd)/terraform:/workspace -v $(pwd)/shared:/shared -e ARM_CLIENT_ID=xxx -e ARM_CLIENT_SECRET=yyy -e ARM_SUBSCRIPTION_ID=zzz -e ARM_TENANT_ID=vvv -w /workspace hashicorp/terraform:0.13.7 plan
docker run --rm -it -v $(pwd)/terraform:/workspace -v $(pwd)/shared:/shared -e ARM_CLIENT_ID=xxx -e ARM_CLIENT_SECRET=yyy -e ARM_SUBSCRIPTION_ID=zzz -e ARM_TENANT_ID=vvv -w /workspace hashicorp/terraform:0.13.7 apply
docker run --rm -it -v $(pwd)/terraform:/workspace -v $(pwd)/shared:/shared -e ARM_CLIENT_ID=xxx -e ARM_CLIENT_SECRET=yyy -e ARM_SUBSCRIPTION_ID=zzz -e ARM_TENANT_ID=vvv -w /workspace hashicorp/terraform:0.13.7 destroy
```
(notice `ARM_` prefix on passed envs)
